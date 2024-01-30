package monitor

import (
	"context"
	"encoding/json"

	"github.com/0xrawsec/golang-etw/etw"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
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

func newSession() (*etw.RealTimeSession, error) {
	s := etw.NewRealTimeSession(etwSessionName)
	if s == nil {
		return nil, errors.New("failed to create ETW session")
	}
	err := enableProvider(s, MicrosoftWindowsKernelProcess, WINEVENT_KEYWORD_PROCESS, WINEVENT_KEYWORD_PROCESS)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func enableProvider(s *etw.RealTimeSession, providerGuid string, matchAnyKeyword, matchAllKeyword uint64) error {
	p := etw.MustParseProvider(providerGuid)
	p.MatchAnyKeyword = matchAnyKeyword
	p.MatchAllKeyword = matchAllKeyword

	log.Infof("Enabling %s provider (GUID: %s, ALL keyword: %x, ANY keyword: %x)", p.Name, p.GUID, p.MatchAllKeyword, p.MatchAnyKeyword)
	err := s.EnableProvider(p)
	if err != nil {
		return errors.Wrap(err, "failed to enable provider")
	}
	log.Infof("Enabled %s provider (GUID: %s)", p.Name, p.GUID)
	return nil
}

func traceAuditEvents(ctx context.Context, ch chan Event) error {
	s, err := newSession()
	if err != nil {
		return err
	}
	defer s.Stop()

	c := etw.NewRealTimeConsumer(ctx).FromSessions(s)
	defer c.Stop()

	go func() {
		for e := range c.Events {
			if ctx.Err() != nil {
				return
			}
			eid := e.System.EventID
			if eid != PROCESS_STARTED && eid != PROCESS_STOPPED {
				continue
			}
			evt, err := parseEvent(e)
			if err != nil {
				log.Errorf("Failed to parse event: %v", err)
				continue
			}
			ch <- *evt
		}
	}()

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

// TODO: process data
func parseEvent(e *etw.Event) (*Event, error) {
	var (
		objectType string
		eventType  string
	)
	if e.System.Task.Name == "ProcessStart" {
		objectType = ObjectTypeProcess
		eventType = EventTypeStarted
	} else if e.System.Task.Name == "ProcessStop" {
		objectType = ObjectTypeProcess
		eventType = EventTypeStopped
	} else {
		return nil, nil
	}
	data, err := marshalEvent(e)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal event")
	}
	pid := int(e.System.Execution.ProcessID)
	evt, err := NewEvent(nil, objectType, eventType, pid, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create event")
	}
	return evt, nil
}

func marshalEvent(e *etw.Event) (map[string]interface{}, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal event")
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal event")
	}
	return m, nil
}
