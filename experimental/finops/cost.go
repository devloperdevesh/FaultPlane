package finops

import (
	"errors"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

// Rates defines explicitly configured runtime cost rates.
//
// CPUHourly is the cost of one fully utilized CPU unit for one hour.
// MemoryGBHourly is the cost of one GB of allocated memory for one hour.
//
// These values are configuration inputs, not fabricated billing data.
type Rates struct {
	CPUHourly      float64
	MemoryGBHourly float64
}

// Snapshot represents a FinOps view derived from actual FaultPlane telemetry.
type Snapshot struct {
	Requests      uint64    `json:"requests"`
	Workers       uint64    `json:"workers"`
	Recoveries    uint64    `json:"recoveries"`
	Checkpoints   uint64    `json:"checkpoints"`
	CPU           float64   `json:"cpu"`
	MemoryMB      float64   `json:"memory_mb"`
	CPUCost       float64   `json:"cpu_cost"`
	MemoryCost    float64   `json:"memory_cost"`
	EstimatedCost float64   `json:"estimated_cost"`
	ObservedAt    time.Time `json:"observed_at"`
}

// Calculate converts one real telemetry snapshot into an estimated runtime cost.
//
// The calculation is intentionally based only on observed runtime telemetry
// and explicitly supplied rates. It does not claim to represent provider
// billing or an actual cloud invoice.
func Calculate(metrics telemetry.RuntimeMetrics, rates Rates) (Snapshot, error) {
	if rates.CPUHourly < 0 {
		return Snapshot{}, errors.New("cpu hourly rate cannot be negative")
	}

	if rates.MemoryGBHourly < 0 {
		return Snapshot{}, errors.New("memory GB hourly rate cannot be negative")
	}

	memoryGB := metrics.Memory / 1024

	cpuCost := metrics.CPU * rates.CPUHourly
	memoryCost := memoryGB * rates.MemoryGBHourly

	return Snapshot{
		Requests:      metrics.Requests,
		Workers:       metrics.Workers,
		Recoveries:    metrics.Recoveries,
		Checkpoints:   metrics.Checkpoints,
		CPU:           metrics.CPU,
		MemoryMB:      metrics.Memory,
		CPUCost:       cpuCost,
		MemoryCost:    memoryCost,
		EstimatedCost: cpuCost + memoryCost,
		ObservedAt:    metrics.UpdatedAt,
	}, nil
}

// FromRegistry creates a FinOps snapshot directly from the live telemetry
// registry.
func FromRegistry(registry *telemetry.Registry, rates Rates) (Snapshot, error) {
	if registry == nil {
		return Snapshot{}, errors.New("telemetry registry is nil")
	}

	return Calculate(registry.Snapshot(), rates)
}
