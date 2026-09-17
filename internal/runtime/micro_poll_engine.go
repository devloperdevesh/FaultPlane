package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// PollFunc performs one bounded polling operation.
// The supplied context is cancelled when the engine stops or the poll timeout expires.
type PollFunc func(context.Context) error

// PollMetrics receives lifecycle metrics from the poll engine.
type PollMetrics interface {
	RecordPoll()
	RecordSuccess()
	RecordError(error)
	RecordDuration(time.Duration)
}

// PollEngineConfig controls the polling lifecycle.
type PollEngineConfig struct {
	Interval    time.Duration
	PollTimeout time.Duration
	Poll        PollFunc
	Metrics     PollMetrics
}

// MicroPollEngine runs bounded, non-overlapping polling work.
type MicroPollEngine struct {
	interval    time.Duration
	pollTimeout time.Duration
	poll        PollFunc
	metrics     PollMetrics

	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan struct{}
	running bool

	polls   atomic.Uint64
	success atomic.Uint64
	errors  atomic.Uint64
}

// NewMicroPollEngine validates configuration and creates an idle engine.
func NewMicroPollEngine(cfg PollEngineConfig) (*MicroPollEngine, error) {
	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("poll interval must be greater than zero")
	}
	if cfg.PollTimeout <= 0 {
		return nil, fmt.Errorf("poll timeout must be greater than zero")
	}
	if cfg.Poll == nil {
		return nil, fmt.Errorf("poll function is required")
	}

	return &MicroPollEngine{
		interval:    cfg.Interval,
		pollTimeout: cfg.PollTimeout,
		poll:        cfg.Poll,
		metrics:     cfg.Metrics,
	}, nil
}

// Start begins the polling lifecycle.
// Start is idempotent while the engine is running.
func (e *MicroPollEngine) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("poll engine start context is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return nil
	}

	runCtx, cancel := context.WithCancel(ctx)
	e.ctx = runCtx
	e.cancel = cancel
	e.done = make(chan struct{})
	e.running = true

	go e.run(runCtx, e.done)

	return nil
}

// Stop cancels the active polling lifecycle and waits for its worker to exit.
// Stop is safe to call repeatedly and concurrently.
func (e *MicroPollEngine) Stop() {
	if e == nil {
		return
	}

	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}

	cancel := e.cancel
	done := e.done
	e.mu.Unlock()

	cancel()
	<-done
}

func (e *MicroPollEngine) run(ctx context.Context, done chan struct{}) {
	defer close(done)

	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.markStopped(done)
			return

		case <-ticker.C:
			e.executePoll(ctx)
		}
	}
}

func (e *MicroPollEngine) markStopped(done chan struct{}) {
	e.mu.Lock()
	if e.done == done {
		e.running = false
		e.ctx = nil
		e.cancel = nil
		e.done = nil
	}
	e.mu.Unlock()
}

func (e *MicroPollEngine) executePoll(parent context.Context) {
	pollCtx, cancel := context.WithTimeout(parent, e.pollTimeout)
	defer cancel()

	start := time.Now()

	e.polls.Add(1)
	if e.metrics != nil {
		e.metrics.RecordPoll()
	}

	err := e.poll(pollCtx)

	duration := time.Since(start)
	if e.metrics != nil {
		e.metrics.RecordDuration(duration)
	}

	if err != nil {
		e.errors.Add(1)
		if e.metrics != nil {
			e.metrics.RecordError(err)
		}
		return
	}

	e.success.Add(1)
	if e.metrics != nil {
		e.metrics.RecordSuccess()
	}
}

// Running reports whether the engine currently has an active lifecycle.
func (e *MicroPollEngine) Running() bool {
	if e == nil {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.running
}

// Stats returns completed poll, successful poll and failed poll counts.
func (e *MicroPollEngine) Stats() (polls, successes, failures uint64) {
	if e == nil {
		return 0, 0, 0
	}

	return e.polls.Load(), e.success.Load(), e.errors.Load()
}
