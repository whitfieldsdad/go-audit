package monitor

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"

	"github.com/Velocidex/etw"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

const (
	etwSessionNamePrefix = "pew-"
)

const (
	WINEVENT_KEYWORD_PROCESS     = 0x10
	WINEVENT_KEYWORD_THREAD      = 0x20
	WINEVENT_KEYWORD_IMAGE       = 0x40
	WindowsKernelProcessAnalytic = 0x8000000000000000
)

const (
	MicrosoftWindowsKernelProcess = "{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}"
)

type WindowsProcessMonitor struct {
	providerGuid windows.GUID
	rawEvents    map[string]chan etw.Event // event type -> buffered channel
}

func NewWindowsProcessMonitor() (*WindowsProcessMonitor, error) {
	providerGuid, err := windows.GUIDFromString(MicrosoftWindowsKernelProcess)
	if err != nil {
		return nil, err
	}
	return &WindowsProcessMonitor{
		providerGuid: providerGuid,
		rawEvents:    make(map[string]chan etw.Event, EventBufferSize),
	}, nil
}

func (m WindowsProcessMonitor) etwCallbackFunction(e *etw.Event) {
	blob, err := json.Marshal(e)
	if err != nil {
		log.Errorf("Failed to JSON marshal event: %v", err)
	}
	log.Infof("Received event: %v", blob)
}

func (m WindowsProcessMonitor) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Create the session.
	session, err := newSession(m.providerGuid, m.etwCallbackFunction)
	if err != nil {
		log.Fatalf("Failed to create ETW session: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	// Start goroutines for:
	// - Reading raw events from the ETW session
	// - Parsing events
	go m.goReadEvents(ctx, session, &wg)
	go m.goParseEvents(ctx, &wg)

	// Cancel the context on SIGINT.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	<-sigCh
	log.Info("Received shutdown signal, shutting down...")
	cancel()

	wg.Wait()
	return nil
}

func (m WindowsProcessMonitor) goReadEvents(ctx context.Context, session *etw.Session, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Info("Processing events...")
		err := session.Process()
		if err != nil {
			log.Fatalf("Failed to process events: %v", err)
		}
		log.Info("Finished processing events")
	}()

	log.Infof("Waiting for shutdown signal...")
	<-ctx.Done()

	log.Infof("Closing session...")
	if err := session.Close(); err != nil {
		log.Fatalf("Failed to close session: %v", err)
	}
	log.Infof("Closed session")
}

func (m WindowsProcessMonitor) goParseEvents(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		for _, c := range m.rawEvents {
			select {
			case e := <-c:
				blob, _ := json.Marshal(e)
				log.Infof("Read audit event: %s", string(blob))
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}
}

type etwCallbackFunction func(e *etw.Event)

func newSession(providerGuid windows.GUID, cb etwCallbackFunction) (*etw.Session, error) {
	var session *etw.Session
	var sessionExistsErr etw.ExistsError
	var err error

	name := getSessionName(providerGuid)
	var createSession = func() (*etw.Session, error) {
		return etw.NewSession(name, printEtwEvent)
	}
	log.Infof("Creating ETW trace session: %s", name)
	session, err = createSession()
	if err != nil {
		if !errors.As(err, &sessionExistsErr) {
			return nil, errors.Wrap(err, "failed to create ETW trace session")
		}
		log.Infof("ETW trace session already exists: %s", name)
		err = etw.KillSession(name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to kill existing ETW trace session")
		}
		log.Infof("Killed existing ETW trace session: %s", name)
		session, err = createSession()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create ETW trace session")
		}
	}
	log.Infof("Created ETW trace session: %s", name)

	// Add the provider.
	log.Infof("Subscribing to ETW provider: %s", providerGuid.String())

	allKeyword := uint64(0x0)
	anyKeyword := uint64(0x50)

	opts := etw.SessionOptions{
		Name:            name,
		Guid:            providerGuid,
		Level:           etw.TRACE_LEVEL_VERBOSE,
		MatchAnyKeyword: anyKeyword,
		MatchAllKeyword: allKeyword,
	}
	err = session.SubscribeToProvider(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to ETW provider")
	}
	log.Infof("Subscribed to ETW provider: %s", providerGuid.String())
	return session, nil
}

func getSessionName(providerGuid windows.GUID) string {
	return etwSessionNamePrefix + providerGuid.String()
}

func printEtwEvent(e *etw.Event) {
	blob, err := json.Marshal(e)
	if err != nil {
		log.Errorf("Failed to JSON marshal event: %v", err)
	}
	log.Infof("Received event: %v", string(blob))
}
