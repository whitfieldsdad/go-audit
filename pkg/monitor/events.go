package monitor

import "time"

type Event struct {
	Header EventHeader `json:"header"`
	Data   EventData   `json:"data"`
}

type EventData interface {
	ToECS() error
}

type EventHeader struct {
	ID         string    `json:"id"`
	Time       time.Time `json:"time"`
	ObjectType string    `json:"object_type"`
	EventType  string    `json:"event_type"`
}

type ProcessStartEventData struct {
	PID        int        `json:"pid"`
	PPID       *int       `json:"ppid,omitempty"`
	CreateTime *time.Time `json:"create_time"`
	Executable *File      `json:"executable,omitempty"`
}

func (e *ProcessStartEventData) ToECS() error {
	panic("implement me")
}

type ProcessStopEventData struct {
	PID        int        `json:"pid"`
	PPID       *int       `json:"ppid,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	ExitTime   *time.Time `json:"exit_time,omitempty"`
	ExitCode   *int       `json:"exit_code,omitempty"`
	Executable *File      `json:"executable,omitempty"`
}

func (e *ProcessStopEventData) ToECS() error {
	panic("implement me")
}
