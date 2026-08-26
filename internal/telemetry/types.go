package telemetry

import "time"

type RuntimeMetrics struct {
	Requests    uint64 `json:"requests"`
	Workers     uint64 `json:"workers"`
	Recoveries  uint64 `json:"recoveries"`
	Checkpoints uint64 `json:"checkpoints"`

	TotalLatencyMs uint64 `json:"total_latency_ms"`
	LatencySamples uint64 `json:"latency_samples"`

	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`

	UpdatedAt time.Time `json:"updated_at"`
}

func (m RuntimeMetrics) AverageLatency() float64 {
	if m.LatencySamples == 0 {
		return 0
	}

	return float64(m.TotalLatencyMs) / float64(m.LatencySamples)
}
