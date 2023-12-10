package monitor

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

var (
	EventBufferSize = 10000
)

type AuditMonitor struct {
	events    chan Event
	eventSink EventSink
}

func NewAuditMonitor(sink EventSink) (*AuditMonitor, error) {
	if _, err := os.Stat(RawAuditEventDir); os.IsNotExist(err) {
		err := os.MkdirAll(RawAuditEventDir, 0755)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create directory for storing raw audit events")
		}
	}
	if sink == nil {
		sink = &StdoutEventSink{}
	}
	return &AuditMonitor{
		eventSink: sink,
		events:    make(chan Event, EventBufferSize),
	}, nil
}

func (m *AuditMonitor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go m.goConsumeRawEvents(ctx, cancel, &wg)
	go m.goProduceRawEvents(ctx, cancel, &wg)

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

func (m *AuditMonitor) goConsumeRawEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	err := readRawAuditEvents(ctx, m.events)
	if err != nil {
		log.Errorf("Failed to read events: %v", err)
		cancel()
	}
}

func (m *AuditMonitor) goProduceRawEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-m.events:
			m.eventSink.Write(ctx, e)
		}
	}
}
