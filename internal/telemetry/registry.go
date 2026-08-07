package telemetry

import (
	"sync"
	"time"
)

// Registry provides concurrent-safe runtime metrics storage.
type Registry struct {
	mu      sync.RWMutex
	metrics RuntimeMetrics
}

// NewRegistry creates a new metrics registry.
func NewRegistry() *Registry {
	return &Registry{
		metrics: RuntimeMetrics{
			UpdatedAt: time.Now(),
		},
	}
}

// IncRequests increments incoming request count.
func (r *Registry) IncRequests() {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Requests++
	r.metrics.UpdatedAt = time.Now()
}

// IncWorkers changes active workers count.
func (r *Registry) IncWorkers() {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Workers++
	r.metrics.UpdatedAt = time.Now()
}

// RecordRecovery records successful recovery.
func (r *Registry) RecordRecovery() {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Recoveries++
	r.metrics.UpdatedAt = time.Now()
}

// RecordCheckpoint records checkpoint creation.
func (r *Registry) RecordCheckpoint() {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Checkpoints++
	r.metrics.UpdatedAt = time.Now()
}

// RecordLatency stores latency sample.
func (r *Registry) RecordLatency(durationMs uint64) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.TotalLatencyMs += durationMs
	r.metrics.LatencySamples++

	r.metrics.UpdatedAt = time.Now()
}

// Snapshot returns current metrics safely.
func (r *Registry) Snapshot() RuntimeMetrics {

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.metrics
}
