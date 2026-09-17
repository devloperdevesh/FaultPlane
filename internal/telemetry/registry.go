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

	kernel KernelMetrics
}

type KernelMetrics struct {
	Transitions       uint64    `json:"transitions"`
	Errors            uint64    `json:"errors"`
	Events            uint64    `json:"events"`
	Enforcement       uint64    `json:"enforcement"`
	Recoveries        uint64    `json:"recoveries"`
	LoaderLoads       uint64    `json:"loader_loads"`
	LoaderCloses      uint64    `json:"loader_closes"`
	LoaderFailures    uint64    `json:"loader_failures"`
	LastState         string    `json:"last_state"`
	LastEvent         string    `json:"last_event"`
	LastReason        string    `json:"last_reason"`
	LastError         string    `json:"last_error"`
	LastTransitionAt  time.Time `json:"last_transition_at"`
	LastEventAt       time.Time `json:"last_event_at"`
	LastErrorAt       time.Time `json:"last_error_at"`
	LastEnforcementAt time.Time `json:"last_enforcement_at"`
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

	r.kernel.Recoveries++
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

func (r *Registry) RecordKernelEvent(eventType string, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.kernel.Events++
	r.kernel.LastEvent = eventType
	r.kernel.LastReason = reason
	r.kernel.LastEventAt = now
	r.metrics.UpdatedAt = now
}

func (r *Registry) RecordKernelTransition(
	from string,
	to string,
	action string,
	eventType string,
	reason string,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.kernel.Transitions++
	r.kernel.LastState = to
	r.kernel.LastEvent = eventType
	r.kernel.LastReason = reason
	r.kernel.LastTransitionAt = now
	r.metrics.UpdatedAt = now

	switch action {
	case "enforce":
		r.kernel.Enforcement++
	case "recover":
		r.kernel.Recoveries++
	}
}

func (r *Registry) RecordKernelError(err error) {
	if err == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.kernel.Errors++
	r.kernel.LastError = err.Error()
	r.kernel.LastErrorAt = now
	r.metrics.UpdatedAt = now
}

func (r *Registry) RecordLoaderLoad() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.kernel.LoaderLoads++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) RecordLoaderClose() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.kernel.LoaderCloses++
	r.metrics.UpdatedAt = time.Now()
}

func (r *Registry) RecordLoaderFailure(err error) {
	if err == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	r.kernel.LoaderFailures++
	r.kernel.Errors++
	r.kernel.LastError = err.Error()
	r.kernel.LastErrorAt = now
	r.metrics.UpdatedAt = now
}

func (r *Registry) KernelSnapshot() KernelMetrics {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.kernel
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
