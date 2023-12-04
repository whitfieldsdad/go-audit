package monitor

import (
	"context"

	"github.com/0xrawsec/golang-etw/etw"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/pkg/errors"
)

type ProcessMonitorInterface interface {
	Run(context.Context) error
}

type ProcessMonitor struct {
	processFilter *ProcessFilter
	events        chan *Event
	rawEvents     chan *etw.Event
	ppidMap       *expirable.LRU[int, int]
}

func NewProcessMonitor(f *ProcessFilter) (*ProcessMonitor, error) {
	ppidMap, err := GetPPIDMap()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get PPID map")
	}
	return &ProcessMonitor{
		processFilter: f,
		events:        make(chan *Event, EventBufferSize),
		rawEvents:     make(chan *etw.Event, EventBufferSize),
		ppidMap:       ppidMap,
	}, nil
}

func (m *ProcessMonitor) AddProcessFilter(filter ProcessFilter) {
	if m.processFilter == nil {
		m.processFilter = &filter
	} else {
		m.processFilter.Merge(filter)
	}
}
