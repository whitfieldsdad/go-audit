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
	ProcessListInterval = 10 * time.Millisecond
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
	ids, err := listProcessIdentities()
	if err != nil {
		log.Errorf("Failed to list processes: %v", err)
		return
	}
	for _, id := range ids {
		seen.Add(id.Hash(), id)
	}

	ticker := time.NewTicker(ProcessListInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ids, err := listProcessIdentities()
			if err != nil {
				log.Errorf("Failed to list processes: %v", err)
				return
			}
			for _, id := range ids {
				h := id.Hash()
				_, ok := seen.Get(h)
				if !ok {
					seen.Add(h, id)

					var process *Process
					process, err = GetProcess(id.PID, opts)
					if err != nil {
						log.Warnf("A new process was detected, but we weren't fast enough to get its details: %v (PID: %d, PPID: %d)", err, id.PID, id.PPID)
						process = &Process{
							PID:  id.PID,
							PPID: id.PPID,
						}
					}
					details := ProcessStartEventData{
						Process: *process,
					}
					log.Infof("Process started (PID: %d, PPID: %d, name: %s)", process.PID, process.PPID, process.Name)
					m.Events <- NewEvent(ObjectTypeProcess, EventTypeStarted, details)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
