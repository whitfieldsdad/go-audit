package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	EventBufferSize = 10000
)

type AuditMonitor struct {
	rawEvents          chan interface{}
	processStartEvents chan *Event
	processStopEvents  chan *Event
	processFilter      *ProcessFilter
}

func NewAuditMonitor(f *ProcessFilter) (*AuditMonitor, error) {
	return &AuditMonitor{
		rawEvents:          make(chan interface{}, EventBufferSize),
		processStartEvents: make(chan *Event, EventBufferSize),
		processStopEvents:  make(chan *Event, EventBufferSize),
		processFilter:      f,
	}, nil
}

func (m *AuditMonitor) AddProcessFilter(filter ProcessFilter) {
	if m.processFilter == nil {
		m.processFilter = &filter
	} else {
		m.processFilter.Merge(filter)
	}
}

func (m *AuditMonitor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(3)

	go m.goReadEvents(ctx, cancel, &wg)
	go m.goParseEvents(ctx, cancel, &wg)
	go m.goProduceEvents(ctx, cancel, &wg)

	// Handle signals.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	log.Info("Waiting for context cancellation or SIGINT...")
	select {
	case <-ctx.Done():
		log.Infof("Context cancelled")
		cancel()
	case <-signalChannel:
		log.Infof("Received SIGINT")
		cancel()
	}
	log.Info("Shutting down...")
	wg.Wait()
	log.Info("Shutdown complete")
	return nil
}

func (m *AuditMonitor) goReadEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	err := readEvents(ctx, m.rawEvents, m.processFilter)
	if err != nil {
		log.Errorf("Failed to read events: %v", err)
		cancel()
	}
}

func (m *AuditMonitor) goParseEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	// Read from the process start/stop channels
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-m.rawEvents:
			if e == nil {
				continue
			}
			event, err := ParseEvent(e)
			if err != nil {
				log.Errorf("Failed to parse event: %v", err)
				continue
			}
			if event == nil {
				continue
			}
			if event.Header.IsProcessStartedEvent() {
				m.processStartEvents <- event
			} else if event.Header.IsProcessStoppedEvent() {
				m.processStopEvents <- event
			} else {
				log.Errorf("Unknown event type: %s %s", event.Header.ObjectType, event.Header.EventType)
			}
		}
	}
}

func (m *AuditMonitor) goProduceEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	var e *Event
	for {
		select {
		case <-ctx.Done():
			return
		case e = <-m.processStartEvents:
			printJSON(e)
		case e = <-m.processStopEvents:
			printJSON(e)
		}
	}
}

func printJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Errorf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(b))
}
