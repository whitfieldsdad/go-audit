package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/0xrawsec/golang-etw/etw"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/whitfieldsdad/go-audit/pkg/util"
)

const (
	etwSessionName = "go-audit"
)

const (
	WINEVENT_KEYWORD_PROCESS     = 0x10
	WINEVENT_KEYWORD_IMAGE       = 0x40
	WindowsKernelProcessAnalytic = 0x8000000000000000
)

const (
	ProcessStart = 1
	ProcessStop  = 2
)

const (
	MicrosoftWindowsKernelProcess = "{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}"
)

var (
	EventBufferSize = 10000
)

type WindowsProcessMonitor struct {
	rawEvents chan *etw.Event
}

func NewWindowsProcessMonitor() (*WindowsProcessMonitor, error) {
	rawEvents := make(chan *etw.Event, EventBufferSize)
	return &WindowsProcessMonitor{
		rawEvents: rawEvents,
	}, nil
}

func (m WindowsProcessMonitor) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Create the ETW session.
	s := etw.NewRealTimeSession(etwSessionName)
	if s == nil {
		return errors.New("failed to create ETW session")
	}
	defer s.Stop()

	// Enable the Microsoft-Windows-Kernel-Process provider.
	p := etw.MustParseProvider(MicrosoftWindowsKernelProcess)
	p.MatchAnyKeyword = WINEVENT_KEYWORD_PROCESS | WINEVENT_KEYWORD_IMAGE

	err := s.EnableProvider(p)
	if err != nil {
		return errors.Wrap(err, "failed to enable Microsoft-Windows-Kernel-Process provider")
	}

	ctx, cancel := context.WithCancel(ctx)

	// Start goroutines for:
	// - Reading raw events from the ETW session
	// - Parsing events
	go m.goReadEvents(ctx, cancel, s, &wg)
	go m.goParseEvents(ctx, &wg)

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

func (m WindowsProcessMonitor) goReadEvents(ctx context.Context, cancel context.CancelFunc, s *etw.RealTimeSession, wg *sync.WaitGroup) {
	defer wg.Done()

	c := etw.NewRealTimeConsumer(ctx)
	defer c.Stop()

	c.FromSessions(s)

	go func() {
		for e := range c.Events {
			t := e.System.EventID
			if !(t == ProcessStart || t == ProcessStop) {
				continue
			}
			m.rawEvents <- e
		}
	}()

	log.Info("Starting consumer...")
	err := c.Start()
	if err != nil {
		log.Errorf("Failed to start consumer: %s", err)
		cancel()
		return
	}
	log.Info("Started consumer")

	<-ctx.Done()
	log.Info("Stopped reading events")
}

func (m WindowsProcessMonitor) goParseEvents(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case r := <-m.rawEvents:
			e, err := m.parseEvent(r)
			if e == nil {
				continue
			}
			if err != nil {
				log.Errorf("Failed to parse event: %s", err)
			}
			b, err := json.Marshal(e)
			if err != nil {
				log.Errorf("Failed to marshal event: %s", err)
			}
			fmt.Println(string(b))
		case <-ctx.Done():
			log.Info("Stopped parsing events")
			return
		}
	}
}

func (m WindowsProcessMonitor) parseEvent(e *etw.Event) (*Event, error) {
	var (
		err        error
		ppid       *int
		createTime *time.Time
	)

	eventId := e.System.EventID
	if !(eventId == ProcessStart || eventId == ProcessStop) {
		return nil, nil
	}
	pid := int(e.System.Execution.ProcessID)

	parentProcessID, ok := e.EventData["ParentProcessID"].(string)
	if ok {
		if ppidInt, err := strconv.Atoi(parentProcessID); err == nil {
			ppid = &ppidInt
		}
	}

	createTimeString, ok := e.EventData["CreateTime"].(string)
	if ok {
		createTime, _ = util.ParseTimestamp(createTimeString)
	}

	executable, _ := GetProcessExecutable(pid)

	if eventId == ProcessStart {
		if ppid == nil {
			ppid, err = GetPPID(pid)
			if err == nil {
				log.Info("Resolved PPID")
			}
		}
		return &Event{
			Header: EventHeader{
				ID:         uuid.New().String(),
				Time:       e.System.TimeCreated.SystemTime,
				ObjectType: "process",
				EventType:  "started",
			},
			Data: &ProcessStartEventData{
				PID:        pid,
				PPID:       ppid,
				CreateTime: createTime,
				Executable: executable,
			},
		}, nil

	} else if eventId == ProcessStop {
		var (
			exitTime *time.Time
			exitCode *int
		)
		exitTimeString, ok := e.EventData["ExitTime"].(string)
		if ok {
			exitTime, _ = util.ParseTimestamp(exitTimeString)
		}
		exitCodeString, ok := e.EventData["ExitStatus"].(string)
		if ok {
			exitCodeInt, _ := strconv.Atoi(exitCodeString)
			if err == nil {
				exitCode = &exitCodeInt
			}
		}
		return &Event{
			Header: EventHeader{
				ID:         uuid.New().String(),
				Time:       e.System.TimeCreated.SystemTime,
				ObjectType: "process",
				EventType:  "stopped",
			},
			Data: &ProcessStopEventData{
				PID:        pid,
				PPID:       ppid,
				CreateTime: createTime,
				ExitTime:   exitTime,
				ExitCode:   exitCode,
				Executable: executable,
			},
		}, nil
	}
	return nil, nil
}
