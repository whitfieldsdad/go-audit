package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/0xrawsec/golang-etw/etw"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/pkg/errors"
	"github.com/whitfieldsdad/go-audit/pkg/util"
)

const (
	etwSessionName = "go-audit"
)

const (
	WINEVENT_KEYWORD_PROCESS                          = 0x10
	WINEVENT_KEYWORD_THREAD                           = 0x20
	WINEVENT_KEYWORD_IMAGE                            = 0x40
	WINEVENT_KEYWORD_CPU_PRIORITY                     = 0x80
	WINEVENT_KEYWORD_OTHER_PRIORITY                   = 0x100
	WINEVENT_KEYWORD_PROCESS_FREEZE                   = 0x200
	WINEVENT_KEYWORD_JOB                              = 0x400
	WINEVENT_KEYWORD_ENABLE_PROCESS_TRACING_CALLBACKS = 0x800
	WINEVENT_KEYWORD_JOB_IO                           = 0x1000
	WINEVENT_KEYWORD_WORK_ON_BEHALF                   = 0x2000
	WINEVENT_KEYWORD_JOB_SILO                         = 0x4000
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

var (
	PidMapTTL     = time.Minute * 15
	PidMapMaxSize = 10000
)

type WindowsProcessMonitor struct {
	processFilter *ProcessFilter
	events        chan *Event
	rawEvents     chan *etw.Event
	ppidMap       *expirable.LRU[int, int]
}

func NewWindowsProcessMonitor(f *ProcessFilter) (*WindowsProcessMonitor, error) {
	ppidMap, err := GetPPIDMap()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get PPID map")
	}
	return &WindowsProcessMonitor{
		processFilter: f,
		events:        make(chan *Event, EventBufferSize),
		rawEvents:     make(chan *etw.Event, EventBufferSize),
		ppidMap:       ppidMap,
	}, nil
}

func (m *WindowsProcessMonitor) AddProcessFilter(filter ProcessFilter) {
	if m.processFilter == nil {
		m.processFilter = &filter
	} else {
		m.processFilter.Merge(filter)
	}
}

func (m *WindowsProcessMonitor) Run(ctx context.Context) error {

	// Create the ETW session.
	s := etw.NewRealTimeSession(etwSessionName)
	if s == nil {
		return errors.New("failed to create ETW session")
	}
	defer s.Stop()

	// Enable the Microsoft-Windows-Kernel-Process provider.
	p := etw.MustParseProvider(MicrosoftWindowsKernelProcess)
	p.MatchAnyKeyword = WINEVENT_KEYWORD_PROCESS | WINEVENT_KEYWORD_IMAGE | WINEVENT_KEYWORD_THREAD

	err := s.EnableProvider(p)
	if err != nil {
		return errors.Wrap(err, "failed to enable Microsoft-Windows-Kernel-Process provider")
	}

	ctx, cancel := context.WithCancel(ctx)

	// Start goroutines for:
	// - Reading events
	// - Parsing events
	// - Emitting events
	var wg sync.WaitGroup
	wg.Add(3)

	go m.goReadEvents(ctx, cancel, s, &wg)
	go m.goParseEvents(ctx, &wg)
	go m.goEmitEvents(ctx, &wg)

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

// TODO
func (m *WindowsProcessMonitor) eventMatchesFilter(e *etw.Event) bool {
	if m.processFilter == nil {
		return true
	}
	entry := m.etwEventToProcessEntry(e)
	if len(m.processFilter.Pids) > 0 && !slices.Contains(m.processFilter.Pids, entry.Pid) {
		return false
	}
	if len(m.processFilter.Ppids) > 0 {
		if entry.PPid == nil || !slices.Contains(m.processFilter.Ppids, *entry.PPid) {
			return false
		}
	}
	return true
}

func (m *WindowsProcessMonitor) etwEventToProcessEntry(e *etw.Event) ProcessEntry {
	pid := int(e.System.Execution.ProcessID)

	var ppid *int
	if e.System.EventID == ProcessStart {
		ppid, _ = extractPPID(e)
	} else {
		ppidptr, ok := m.ppidMap.Get(pid)
		if ok {
			ppid = &ppidptr
			log.Infof("Resolved PPID from PPID map (PID: %d, PPID: %d)", pid, ppid)
		}
	}
	return ProcessEntry{
		Pid:  pid,
		PPid: ppid,
	}
}

// TODO
func (m *WindowsProcessMonitor) goReadEvents(ctx context.Context, cancel context.CancelFunc, s *etw.RealTimeSession, wg *sync.WaitGroup) {
	defer wg.Done()

	c := etw.NewRealTimeConsumer(ctx)
	defer c.Stop()

	c.FromSessions(s)

	go func() {
		for e := range c.Events {
			pid := int(e.System.Execution.ProcessID)
			eventId := e.System.EventID
			if !(eventId == ProcessStart || eventId == ProcessStop) {
				continue
			}
			log.Debugf("Read event (ID: %d, keyword name: %s, task name: %s, PID: %d)", eventId, e.System.Keywords.Name, e.System.Task.Name, pid)

			// Not all events contain a PPID, so we need to keep track of PID -> PPID mappings ourselves.
			m.updatePPIDMap(e)

			if !m.eventMatchesFilter(e) {
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

func (m *WindowsProcessMonitor) updatePPIDMap(e *etw.Event) {
	pid := int(e.System.Execution.ProcessID)
	eid := e.System.EventID
	if eid == ProcessStop {
		log.Infof("Removing PID from PPID map (PID: %d, event ID: %d, size: %d)", pid, eid, m.ppidMap.Len())
		m.ppidMap.Remove(pid)
		return
	}
	var ppid *int
	if e.System.EventID == ProcessStart {
		ppidstr, ok := e.EventData["ParentProcessID"].(string)
		if ok {
			ppidptr, err := strconv.Atoi(ppidstr)
			if err == nil {
				ppid = &ppidptr
			}
		}
	}
	if ppid != nil {
		pid := int(e.System.Execution.ProcessID)
		log.Infof("Adding PID to PPID map (PID: %d, PPID: %d, event ID: %d, size: %d)", pid, ppid, e.System.EventID, m.ppidMap.Len())
		m.ppidMap.Add(pid, *ppid)
	}
}

func (m *WindowsProcessMonitor) goParseEvents(ctx context.Context, wg *sync.WaitGroup) {
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
			m.events <- e
		case <-ctx.Done():
			log.Info("Stopped parsing events")
			return
		}
	}
}

func (m *WindowsProcessMonitor) parseEvent(e *etw.Event) (*Event, error) {
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

func (m *WindowsProcessMonitor) goEmitEvents(ctx context.Context, wg *sync.WaitGroup) {
	var (
		e   *Event
		err error
		b   []byte
	)
	defer wg.Done()

	for {
		select {
		case e = <-m.events:
			b, err = json.Marshal(e)
			if err != nil {
				log.Errorf("Failed to marshal event: %s", err)
			}
			fmt.Println(string(b))
		case <-ctx.Done():
			log.Info("Stopped emitting events")
			return
		}
	}
}

func extractPPID(e *etw.Event) (*int, error) {
	ppidstr, ok := e.EventData["ParentProcessID"].(string)
	if !ok {
		return nil, nil
	}
	ppid, err := strconv.Atoi(ppidstr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PPID")
	}
	return &ppid, nil
}
