package monitor

import (
	"time"

	"github.com/whitfieldsdad/go-building-blocks/pkg/bb"
)

type ObjectType int

const (
	ObjectTypeProcess ObjectType = iota + 1
)

func (t ObjectType) String() string {
	switch t {
	case ObjectTypeProcess:
		return "process"
	default:
		return "unknown"
	}
}

type EventType int

const (
	EventTypeStarted EventType = iota + 1
	EventTypeStopped
)

func (t EventType) String() string {
	switch t {
	case EventTypeStarted:
		return "started"
	case EventTypeStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type Event struct {
	Header EventHeader `json:"header"`
	Data   EventData   `json:"data"`
}

type EventHeader struct {
	Id         string    `json:"id"`
	Time       time.Time `json:"time"`
	ObjectType string    `json:"object_type"`
	EventType  string    `json:"event_type"`
}

func (h *EventHeader) IsProcessEvent() bool {
	return h.ObjectType == ObjectTypeProcess.String()
}

func (h *EventHeader) IsProcessStartedEvent() bool {
	return h.IsProcessEvent() && h.EventType == EventTypeStarted.String()
}

func (h *EventHeader) IsProcessStoppedEvent() bool {
	return h.IsProcessEvent() && h.EventType == EventTypeStopped.String()
}

type EventData struct {
	Process *bb.Process `json:"process"`
}
