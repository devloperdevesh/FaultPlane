package finops

import (
	"testing"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func TestCalculateUsesRealTelemetry(t *testing.T) {
	metrics := telemetry.RuntimeMetrics{
		Requests:    12,
		Workers:     3,
		Recoveries:  2,
		Checkpoints: 4,
		CPU:         0.5,
		Memory:      1024,
		UpdatedAt:   time.Unix(100, 0).UTC(),
	}

	snapshot, err := Calculate(metrics, Rates{
		CPUHourly:      2,
		MemoryGBHourly: 1,
	})
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}

	if snapshot.Requests != 12 {
		t.Fatalf("requests = %d, want 12", snapshot.Requests)
	}

	if snapshot.Workers != 3 {
		t.Fatalf("workers = %d, want 3", snapshot.Workers)
	}

	if snapshot.CPUCost != 1 {
		t.Fatalf("cpu cost = %v, want 1", snapshot.CPUCost)
	}

	if snapshot.MemoryCost != 1 {
		t.Fatalf("memory cost = %v, want 1", snapshot.MemoryCost)
	}

	if snapshot.EstimatedCost != 2 {
		t.Fatalf("estimated cost = %v, want 2", snapshot.EstimatedCost)
	}

	if !snapshot.ObservedAt.Equal(metrics.UpdatedAt) {
		t.Fatalf("observed timestamp was not preserved")
	}
}

func TestCalculateRejectsNegativeRates(t *testing.T) {
	metrics := telemetry.RuntimeMetrics{}

	if _, err := Calculate(metrics, Rates{CPUHourly: -1}); err == nil {
		t.Fatal("expected negative CPU rate to fail")
	}

	if _, err := Calculate(metrics, Rates{MemoryGBHourly: -1}); err == nil {
		t.Fatal("expected negative memory rate to fail")
	}
}

func TestFromRegistryUsesLiveRegistry(t *testing.T) {
	registry := telemetry.NewRegistry()

	registry.IncRequests()
	registry.IncWorkers()
	registry.RecordCheckpoint()
	registry.RecordRecovery()

	snapshot, err := FromRegistry(registry, Rates{
		CPUHourly:      1,
		MemoryGBHourly: 1,
	})
	if err != nil {
		t.Fatalf("from registry: %v", err)
	}

	if snapshot.Requests != 1 {
		t.Fatalf("requests = %d, want 1", snapshot.Requests)
	}

	if snapshot.Workers != 1 {
		t.Fatalf("workers = %d, want 1", snapshot.Workers)
	}

	if snapshot.Checkpoints != 1 {
		t.Fatalf("checkpoints = %d, want 1", snapshot.Checkpoints)
	}

	if snapshot.Recoveries != 1 {
		t.Fatalf("recoveries = %d, want 1", snapshot.Recoveries)
	}
}
