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
	rawEvents chan map[string]interface{}
}

func NewAuditMonitor() (*AuditMonitor, error) {
	return &AuditMonitor{
		rawEvents: make(chan map[string]interface{}, EventBufferSize),
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

	err := readRawEvents(ctx, m.rawEvents)
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
		case e := <-m.rawEvents:
			b, _ := json.Marshal(e)
			fmt.Println(string(b))
		}
	}
}
