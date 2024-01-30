package monitor

import (
	"time"

	"github.com/google/uuid"
)

type ObjectType string

const (
	ObjectTypeProcess = "process"
)

type EventType string

const (
	EventTypeStarted = "started"
	EventTypeStopped = "stopped"
)

type Event struct {
	Header EventHeader `json:"header"`
	Data   interface{} `json:"data"`
}

type ProcessStartEventData struct {
	Process
}

type ProcessStopEventData struct {
	PID        int32      `json:"pid"`
	PPID       *int32     `json:"ppid,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	ExitTime   *time.Time `json:"exit_time,omitempty"`
}

type EventHeader struct {
	Id         string     `json:"id"`
	Time       time.Time  `json:"time"`
	ObjectType ObjectType `json:"object_type"`
	EventType  EventType  `json:"event_type"`
}

func NewEvent(objectType ObjectType, eventType EventType, details interface{}) Event {
	return Event{
		Header: EventHeader{
			Id:         uuid.New().String(),
			Time:       time.Now(),
			ObjectType: objectType,
			EventType:  eventType,
		},
		Data: details,
	}
}
