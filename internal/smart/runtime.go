package smart

import (
	"context"
	"sync"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Runtime struct {
	registry *telemetry.Registry
	detector *Detector
	interval time.Duration

	mu     sync.RWMutex
	latest Result
}

func NewRuntime(
	registry *telemetry.Registry,
	detector *Detector,
	interval time.Duration,
) *Runtime {
	if interval <= 0 {
		interval = time.Second
	}

	return &Runtime{
		registry: registry,
		detector: detector,
		interval: interval,
	}
}

func (r *Runtime) Start(ctx context.Context) {
	if r == nil || r.registry == nil || r.detector == nil {
		return
	}

	r.observe()

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.observe()
		}
	}
}

func (r *Runtime) observe() {
	metrics := r.registry.Snapshot()

	snapshot := Snapshot{
		Timestamp: metrics.UpdatedAt,
		Requests:  float64(metrics.Requests),
		Latency:   metrics.AverageLatency(),
		CPU:       metrics.CPU,
		Memory:    metrics.Memory,
	}

	result := r.detector.Observe(snapshot)

	r.mu.Lock()
	r.latest = result
	r.mu.Unlock()
}

func (r *Runtime) Latest() Result {
	if r == nil {
		return Result{}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.latest
}
