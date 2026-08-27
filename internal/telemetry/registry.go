package telemetry

import (
	"runtime"
	"sync"
	"time"
)

type RuntimeCPUSampler interface {
	Sample() float64
}

type Registry struct {
	mu      sync.RWMutex
	metrics RuntimeMetrics
	cpu     RuntimeCPUSampler
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: RuntimeMetrics{
			UpdatedAt: time.Now(),
		},
	}
}

func (r *Registry) SetCPUSampler(sampler RuntimeCPUSampler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cpu = sampler
}

func (r *Registry) IncRequests() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Requests++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) IncWorkers() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Workers++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) RecordRecovery() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Recoveries++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) RecordCheckpoint() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Checkpoints++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) RecordLatency(durationMs uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.TotalLatencyMs += durationMs
	r.metrics.LatencySamples++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) RefreshRuntimeMetrics() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	r.mu.Lock()
	sampler := r.cpu
	r.mu.Unlock()

	var cpu float64
	if sampler != nil {
		cpu = sampler.Sample()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Memory = float64(mem.Alloc) / 1024 / 1024
	r.metrics.CPU = cpu
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) Snapshot() RuntimeMetrics {
	r.RefreshRuntimeMetrics()

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.metrics
}
