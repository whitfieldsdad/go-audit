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
	etwSessionName = "go-audit"
)

const (
	WINEVENT_KEYWORD_PROCESS     = 0x10
	WINEVENT_KEYWORD_THREAD      = 0x20
	WINEVENT_KEYWORD_IMAGE       = 0x40
	WindowsKernelProcessAnalytic = 0x8000000000000000
)

const (
	ProcessStart = 1
)

const (
	MicrosoftWindowsKernelProcess = "{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}"
)

type WindowsProcessMonitor struct {
	providerGuid windows.GUID
	rawEvents    map[uint16]chan etw.Event // event type -> buffered channel
}

func NewWindowsProcessMonitor() (*WindowsProcessMonitor, error) {
	providerGuid, err := windows.GUIDFromString(MicrosoftWindowsKernelProcess)
	if err != nil {
		return nil, err
	}
	return &WindowsProcessMonitor{
		providerGuid: providerGuid,
		rawEvents:    make(map[uint16]chan etw.Event, EventBufferSize),
	}, nil
}

func (m WindowsProcessMonitor) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)

	// Create the session.
	session, err := newSession(m.providerGuid, m.processEvent)
	if err != nil {
		log.Fatalf("Failed to create ETW session: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	// Start goroutines for:
	// - Reading raw events from the ETW session
	// - Parsing events
	go m.goReadEvents(ctx, cancel, session, &wg)
	//go m.goParseEvents(ctx, &wg)

	// Cancel the context on SIGINT.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Check if either the context was cancelled, or we received a SIGINT.
	log.Info("Waiting for context cancellation or SIGINT...")
	select {
	case <-ctx.Done():
		log.Infof("Context cancelled")
		cancel()
	case <-sigCh:
		log.Infof("Received SIGINT")
		cancel()
	}
	log.Info("Shutting down...")
	wg.Wait()
	log.Info("Shutdown complete")
	return nil
}

func (m WindowsProcessMonitor) goReadEvents(ctx context.Context, cancel context.CancelFunc, session *etw.Session, wg *sync.WaitGroup) {
	defer wg.Done()

	killed := false

	go func() {
		log.Info("Reading events...")
		err := session.Process()
		if err != nil {
			log.Fatalf("Failed to process events: %v", err)
		}
		killed = true
		cancel()
	}()

	<-ctx.Done()

	if killed {
		log.Warn("session killed by external process")
	} else {
		log.Info("Closing session...")

		// TODO: remove session.UnsubscribeFromProvider() after recursive locking bug is fixed in session.Close() function
		err := session.UnsubscribeFromProvider(m.providerGuid)
		if err != nil {
			log.Errorf("Failed to unsubscribe from provider: %v", err)
		}
		err = session.Close()
		if err != nil {
			log.Errorf("Failed to close session: %v", err)
		}
	}
	log.Info("Stopped reading events")
}

func (m WindowsProcessMonitor) processEvent(e *etw.Event) {
	t := e.Header.EventDescriptor.ID
	if t != 1 {
		return
	}
	printEtwEvent(e)
	if _, ok := m.rawEvents[t]; !ok {
		m.rawEvents[t] = make(chan etw.Event, EventBufferSize)
	}
	m.rawEvents[t] <- *e
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

func newSession(providerGuid windows.GUID, cb etw.EventCallback) (*etw.Session, error) {
	var session *etw.Session
	var sessionExistsErr etw.ExistsError
	var err error

	var createSession = func() (*etw.Session, error) {
		return etw.NewSession(etwSessionName, cb)
	}
	log.Infof("Creating session: %s", etwSessionName)
	session, err = createSession()
	if err != nil {
		if !errors.As(err, &sessionExistsErr) {
			return nil, errors.Wrap(err, "failed to create session")
		}
		log.Infof("session already exists: %s", etwSessionName)
		err = etw.KillSession(etwSessionName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to kill existing session")
		}
		log.Infof("Killed existing session: %s", etwSessionName)
		session, err = createSession()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create session")
		}
	}
	log.Infof("Created session: %s", etwSessionName)

	// Add the provider.
	log.Infof("Subscribing to ETW provider: %s", providerGuid.String())

	allKeyword := uint64(0x0)
	anyKeyword := uint64(0x50)

	opts := etw.SessionOptions{
		Name:            etwSessionName,
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

func printEtwEvent(e *etw.Event) {
	data, err := e.EventProperties(false)
	if err != nil {
		log.Errorf("Failed to get event properties: %v", err)
		return
	}
	event := map[string]interface{}{
		"Header": e.Header,
		"Data":   data,
	}
	blob, _ := json.Marshal(event)
	log.Infof("Read audit event: %s", string(blob))
}
