package monitor

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	lru "github.com/hashicorp/golang-lru"
)

var (
	ProcessListInterval = 100 * time.Millisecond
	EventBufferSize     = 10000
)

type AuditMonitor struct {
	Events chan Event
}

func NewAuditMonitor() (*AuditMonitor, error) {
	return &AuditMonitor{
		Events: make(chan Event, EventBufferSize),
	}, nil
}

func (m *AuditMonitor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

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
	wg.Wait()
	return nil
}

func (m *AuditMonitor) goReadEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	err := traceAuditEvents(ctx, m.Events)
	if err != nil {
		log.Warnf("Failed to trace audit events: %s", err)
		log.Info("Falling back to polling audit events")
		m.pollAuditEvents(ctx, cancel, wg)
	}
}

func (m *AuditMonitor) pollAuditEvents(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	seen, err := lru.New(10000)
	if err != nil {
		log.Errorf("Failed to create LRU cache: %v", err)
		return
	}
	opts := &ProcessOptions{
		IncludeHashes: false,
	}
	ps, err := ListProcesses(opts)
	if err != nil {
		log.Errorf("Failed to list processes: %v", err)
		return
	}
	for _, p := range ps {
		seen.Add(p.Hash(), p)
	}

	ticker := time.NewTicker(ProcessListInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ps, err := ListProcesses(opts)
			if err != nil {
				log.Errorf("Failed to list processes: %v", err)
				return
			}
			for _, p := range ps {
				_, ok := seen.Get(p.Hash())
				if !ok {
					seen.Add(p.Hash(), p)
					d := ProcessStartEventData{
						PID:        p.PID,
						PPID:       p.PPID,
						Name:       p.Name,
						CreateTime: p.CreateTime,
						Executable: p.Executable,
					}
					log.Infof("Process started (PID: %d, PPID: %d, name: %s)", d.PID, d.PPID, d.Name)

					m.Events <- NewEvent(ObjectTypeProcess, EventTypeStarted, d)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
