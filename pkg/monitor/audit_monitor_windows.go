package monitor

import (
	"context"
	"strconv"

	"github.com/0xrawsec/golang-etw/etw"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/whitfieldsdad/go-audit/pkg/util"
)

const (
	etwSessionName = "go-audit"
)

const (
	WINEVENT_KEYWORD_PROCESS = 0x10
)

const (
	PROCESS_STARTED = 1
	PROCESS_STOPPED = 2
)

const (
	MicrosoftWindowsKernelProcess = "{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}"
)

func readEvents(ctx context.Context, ch chan interface{}, f *ProcessFilter) error {
	s := etw.NewRealTimeSession(etwSessionName)
	if s == nil {
		return errors.New("failed to create ETW session")
	}
	defer s.Stop()

	p := etw.MustParseProvider(MicrosoftWindowsKernelProcess)
	p.MatchAnyKeyword = WINEVENT_KEYWORD_PROCESS

	err := s.EnableProvider(p)
	if err != nil {
		return errors.Wrap(err, "failed to enable Microsoft-Windows-Kernel-Process provider")
	}
	c := etw.NewRealTimeConsumer(ctx).FromSessions(s)
	defer c.Stop()

	go func() {
		for e := range c.Events {
			if ctx.Err() != nil {
				return
			}
			if f != nil {
				pid := int(e.System.Execution.ProcessID)
				ppid, _ := extractPPID(e)
				ref := Process{
					Pid:  pid,
					Ppid: ppid,
				}
				ok, err := f.Matches(ref, nil)
				if err != nil {
					continue
				}
				if !ok {
					continue
				}
			}
			eid := e.System.EventID
			if eid == PROCESS_STARTED || eid == PROCESS_STOPPED {
				ch <- e
			}
		}
	}()

	// Start the consumer.
	log.Infof("Starting ETW consumer...")
	err = c.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start ETW consumer")
	}
	log.Infof("Reading events...")
	<-ctx.Done()
	log.Infof("Stopped reading events")
	return nil
}

func ParseEvent(e interface{}) (*Event, error) {
	return parseEvent(e.(*etw.Event))
}

func parseEvent(e *etw.Event) (*Event, error) {
	process, err := extractProcess(e)
	if process == nil || err != nil {
		return nil, err
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	var (
		objectType ObjectType
		eventType  EventType
	)
	if e.System.EventID == PROCESS_STARTED {
		objectType = ObjectTypeProcess
		eventType = EventTypeStarted
	} else if e.System.EventID == PROCESS_STOPPED {
		objectType = ObjectTypeProcess
		eventType = EventTypeStopped
	} else {
		return nil, nil
	}
	return &Event{
		Header: EventHeader{
			Id:         id.String(),
			Time:       e.System.TimeCreated.SystemTime,
			EventType:  eventType.String(),
			ObjectType: objectType.String(),
		},
		Data: EventData{
			Process: process,
		},
	}, nil
}

func extractProcess(e *etw.Event) (*Process, error) {
	var (
		p   *Process
		err error
	)
	pid := int(e.System.Execution.ProcessID)
	eid := e.System.EventID
	if eid == PROCESS_STOPPED {
		p, err = extractProcessFromProcessStopEvent(e)
		if err != nil {
			return nil, errors.Wrap(err, "failed to extract process from process stop event")
		}
	} else {
		p, err = GetProcess(pid)
		if err != nil && err != process.ErrorProcessNotRunning {
			return nil, err
		}
	}
	return p, nil
}

func extractProcessFromProcessStopEvent(e *etw.Event) (*Process, error) {
	pid := int(e.System.Execution.ProcessID)
	createTime, err := util.TimeFromRFC3339(e.EventData["CreateTime"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse create time")
	}
	exitTime, err := util.TimeFromRFC3339(e.EventData["ExitTime"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse exit time")
	}
	var exitCode int
	var exitCodePtr *int
	s, ok := e.EventData["ExitStatus"].(string)
	if ok {
		exitCode, err = strconv.Atoi(s)
		if err == nil {
			exitCodePtr = &exitCode
		}
	}
	return &Process{
		Pid:        pid,
		CreateTime: createTime,
		ExitTime:   exitTime,
		ExitCode:   exitCodePtr,
	}, nil
}

func extractPPID(e *etw.Event) (*int, error) {
	s, ok := e.EventData["ParentProcessID"].(string)
	if !ok {
		return nil, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PPID")
	}
	return &v, nil
}
