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
	Events chan Event
}

func NewAuditMonitor() (*AuditMonitor, error) {
	if _, err := os.Stat(RawAuditEventDir); os.IsNotExist(err) {
		err := os.MkdirAll(RawAuditEventDir, 0755)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create directory for storing raw audit events")
		}
	}
	return &AuditMonitor{
		Events: make(chan Event, EventBufferSize),
	}, nil
}

func (m *AuditMonitor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go m.goReadEvents(ctx, cancel, &wg)

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

	err := readRawAuditEvents(ctx, m.events)
	if err != nil {
		log.Errorf("Failed to read events: %v", err)
		cancel()
	}
}
