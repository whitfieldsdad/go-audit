package monitor

import (
	"time"

	"github.com/whitfieldsdad/go-building-blocks/pkg/bb"
)

const (
	ObjectTypeProcess = "process"
)

const (
	EventTypeStarted = "started"
	EventTypeStopped = "stopped"
)

type Event struct {
	Header EventHeader            `json:"header"`
	Data   map[string]interface{} `json:"data"`
}

type EventHeader struct {
	Id         string    `json:"id"`
	Time       time.Time `json:"time"`
	Pid        int       `json:"pid"`
	ObjectType string    `json:"object_type"`
	EventType  string    `json:"event_type"`
}

func NewEvent(timestampptr *time.Time, objectType, eventType string, pid int, data map[string]interface{}) (*Event, error) {
	var timestamp time.Time
	if timestampptr == nil {
		timestamp = time.Now()
	} else {
		timestamp = *timestampptr
	}
	header := EventHeader{
		Id:         bb.NewUUID4(),
		Time:       timestamp,
		Pid:        pid,
		ObjectType: objectType,
		EventType:  eventType,
	}
	return &Event{
		Header: header,
		Data:   data,
	}, nil
}

func NewProcessStartEvent(timestampptr *time.Time, pid int, data map[string]interface{}) (*Event, error) {
	return NewEvent(timestampptr, ObjectTypeProcess, EventTypeStarted, pid, data)
}

func NewProcessStopEvent(timestampptr *time.Time, pid int, data map[string]interface{}) (*Event, error) {
	return NewEvent(timestampptr, ObjectTypeProcess, EventTypeStopped, pid, data)
}
