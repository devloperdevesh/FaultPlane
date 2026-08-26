package telemetry

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Event represents one runtime telemetry sample.
type Event struct {
	Name      string
	Timestamp time.Time
	Value     float64
	Metadata  map[string]string
}

// Pipeline asynchronously processes runtime telemetry events.
type Pipeline struct {
	registry *Registry
	events   chan Event

	wg sync.WaitGroup

	mu      sync.RWMutex
	running bool
}

// NewPipeline creates a telemetry pipeline.
func NewPipeline(
	registry *Registry,
	bufferSize int,
) *Pipeline {
	if bufferSize <= 0 {
		bufferSize = 128
	}

	return &Pipeline{
		registry: registry,
		events:   make(chan Event, bufferSize),
	}
}

// Start begins consuming telemetry events.
func (p *Pipeline) Start(ctx context.Context) {
	p.mu.Lock()

	if p.running {
		p.mu.Unlock()
		return
	}

	p.running = true
	p.wg.Add(1)

	p.mu.Unlock()

	go func() {
		defer p.wg.Done()

		for {
			select {
			case <-ctx.Done():
				p.mu.Lock()
				p.running = false
				p.mu.Unlock()
				return

			case event := <-p.events:
				p.process(event)
			}
		}
	}()
}

// Publish adds an event to the pipeline.
func (p *Pipeline) Publish(event Event) error {
	p.mu.RLock()
	running := p.running
	p.mu.RUnlock()

	if !running {
		return errors.New("telemetry pipeline is not running")
	}

	select {
	case p.events <- event:
		return nil
	default:
		return errors.New("telemetry pipeline buffer is full")
	}
}

// Wait waits until the pipeline worker exits.
func (p *Pipeline) Wait() {
	p.wg.Wait()
}

// process translates supported events into registry metrics.
func (p *Pipeline) process(event Event) {
	if p.registry == nil {
		return
	}

	switch event.Name {
	case "request":
		p.registry.IncRequests()

	case "checkpoint":
		p.registry.RecordCheckpoint()

	case "recovery":
		p.registry.RecordRecovery()

	case "latency":
		if event.Value >= 0 {
			p.registry.RecordLatency(uint64(event.Value))
		}
	}
}
