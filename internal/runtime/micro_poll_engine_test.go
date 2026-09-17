package runtime

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testPollMetrics struct {
	polls     atomic.Uint64
	successes atomic.Uint64
	errors    atomic.Uint64
	duration  atomic.Uint64
}

func (m *testPollMetrics) RecordPoll() {
	m.polls.Add(1)
}

func (m *testPollMetrics) RecordSuccess() {
	m.successes.Add(1)
}

func (m *testPollMetrics) RecordError(error) {
	m.errors.Add(1)
}

func (m *testPollMetrics) RecordDuration(d time.Duration) {
	m.duration.Add(uint64(d))
}

func TestMicroPollEngineValidation(t *testing.T) {
	poll := func(context.Context) error { return nil }

	tests := []PollEngineConfig{
		{Interval: 0, PollTimeout: time.Second, Poll: poll},
		{Interval: time.Second, PollTimeout: 0, Poll: poll},
		{Interval: time.Second, PollTimeout: time.Second},
	}

	for _, cfg := range tests {
		if _, err := NewMicroPollEngine(cfg); err == nil {
			t.Fatal("expected configuration error")
		}
	}
}

func TestMicroPollEngineLifecycle(t *testing.T) {
	var calls atomic.Uint64

	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    10 * time.Millisecond,
		PollTimeout: 100 * time.Millisecond,
		Poll: func(ctx context.Context) error {
			calls.Add(1)
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := engine.Start(ctx); err != nil {
		t.Fatal(err)
	}

	if err := engine.Start(ctx); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for calls.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	if calls.Load() == 0 {
		t.Fatal("poll did not execute")
	}

	engine.Stop()

	if engine.Running() {
		t.Fatal("engine still running after Stop")
	}

	before := calls.Load()
	time.Sleep(30 * time.Millisecond)

	if calls.Load() != before {
		t.Fatal("polling continued after Stop")
	}
}

func TestMicroPollEngineContextCancellation(t *testing.T) {
	started := make(chan struct{})

	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    5 * time.Millisecond,
		PollTimeout: time.Second,
		Poll: func(ctx context.Context) error {
			select {
			case <-started:
			default:
				close(started)
			}

			<-ctx.Done()
			return ctx.Err()
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := engine.Start(ctx); err != nil {
		t.Fatal(err)
	}

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("poll never started")
	}

	cancel()
	engine.Stop()

	if engine.Running() {
		t.Fatal("engine still running")
	}
}

func TestMicroPollEnginePollTimeout(t *testing.T) {
	var seen atomic.Bool

	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    5 * time.Millisecond,
		PollTimeout: 20 * time.Millisecond,
		Poll: func(ctx context.Context) error {
			seen.Store(true)
			<-ctx.Done()
			return ctx.Err()
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := engine.Start(ctx); err != nil {
		t.Fatal(err)
	}

	time.Sleep(80 * time.Millisecond)
	engine.Stop()

	if !seen.Load() {
		t.Fatal("poll never executed")
	}

	_, _, failures := engine.Stats()
	if failures == 0 {
		t.Fatal("expected timeout to be recorded as failure")
	}
}

func TestMicroPollEngineErrorAndMetrics(t *testing.T) {
	metrics := &testPollMetrics{}
	expected := errors.New("poll failed")

	var calls atomic.Uint64

	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    5 * time.Millisecond,
		PollTimeout: 100 * time.Millisecond,
		Poll: func(context.Context) error {
			time.Sleep(1 * time.Millisecond)

			if calls.Add(1)%2 == 0 {
				return expected
			}

			return nil
		},
		Metrics: metrics,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := engine.Start(ctx); err != nil {
		t.Fatal(err)
	}

	time.Sleep(70 * time.Millisecond)
	engine.Stop()

	polls, successes, failures := engine.Stats()

	if polls == 0 {
		t.Fatal("expected polls")
	}
	if successes == 0 {
		t.Fatal("expected successful polls")
	}
	if failures == 0 {
		t.Fatal("expected failed polls")
	}
	if metrics.polls.Load() != polls {
		t.Fatal("poll metric mismatch")
	}
	if metrics.successes.Load() != successes {
		t.Fatal("success metric mismatch")
	}
	if metrics.errors.Load() != failures {
		t.Fatal("error metric mismatch")
	}
	if metrics.duration.Load() == 0 {
		t.Fatal("expected duration metric")
	}
}
func TestMicroPollEngineNoOverlappingPolls(t *testing.T) {
	var active atomic.Int32
	var maxActive atomic.Int32

	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    5 * time.Millisecond,
		PollTimeout: 100 * time.Millisecond,
		Poll: func(ctx context.Context) error {
			current := active.Add(1)
			for {
				old := maxActive.Load()
				if current <= old || maxActive.CompareAndSwap(old, current) {
					break
				}
			}

			time.Sleep(20 * time.Millisecond)
			active.Add(-1)
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := engine.Start(ctx); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)
	engine.Stop()

	if maxActive.Load() > 1 {
		t.Fatalf("polls overlapped: max active=%d", maxActive.Load())
	}
}

func TestMicroPollEngineConcurrentStop(t *testing.T) {
	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    5 * time.Millisecond,
		PollTimeout: 50 * time.Millisecond,
		Poll: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Millisecond):
				return nil
			}
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := engine.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			engine.Stop()
		}()
	}

	wg.Wait()

	if engine.Running() {
		t.Fatal("engine still running")
	}
}

func TestMicroPollEngineAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	engine, err := NewMicroPollEngine(PollEngineConfig{
		Interval:    time.Second,
		PollTimeout: time.Second,
		Poll:        func(context.Context) error { return nil },
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := engine.Start(ctx); err == nil {
		t.Fatal("expected cancelled-context error")
	}
}
