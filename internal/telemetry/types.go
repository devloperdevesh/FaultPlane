package telemetry

import "time"

// RuntimeMetrics represents runtime health
// information collected from AgentMesh components.
//
// These metrics are transport independent.
// API, Prometheus and Dashboard layers
// can consume this model later.
type RuntimeMetrics struct {

	// Total incoming requests handled
	Requests uint64 `json:"requests"`

	// Currently active workers
	Workers uint64 `json:"workers"`

	// Successful recovery operations
	Recoveries uint64 `json:"recoveries"`

	// Total checkpoints created
	Checkpoints uint64 `json:"checkpoints"`

	// Total request latency in milliseconds
	TotalLatencyMs uint64 `json:"total_latency_ms"`

	// Number of latency samples
	LatencySamples uint64 `json:"latency_samples"`

	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
}

// AverageLatency calculates average runtime latency.
func (m RuntimeMetrics) AverageLatency() float64 {

	if m.LatencySamples == 0 {
		return 0
	}

	return float64(m.TotalLatencyMs) /
		float64(m.LatencySamples)
}
