package telemetry

import (
	"sync"
	"testing"
	"time"
)

func TestPerformanceMetricsRecordLatency(t *testing.T) {
	registry := NewRegistry()

	registry.RecordLatency(2)
	registry.RecordLatency(4)
	registry.RecordLatency(6)

	snapshot := registry.Snapshot()

	if snapshot.LatencySamples != 3 {
		t.Fatalf("expected 3 latency samples, got %d", snapshot.LatencySamples)
	}

	if snapshot.TotalLatencyMs != 12 {
		t.Fatalf("expected 12ms total latency, got %d", snapshot.TotalLatencyMs)
	}

	if got := snapshot.AverageLatency(); got != 4 {
		t.Fatalf("expected 4ms average latency, got %v", got)
	}
}

func TestPerformanceMetricsConcurrentRecording(t *testing.T) {
	registry := NewRegistry()

	const workers = 8
	const samplesPerWorker = 1000

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < samplesPerWorker; j++ {
				registry.RecordLatency(1)
			}
		}()
	}

	wg.Wait()

	snapshot := registry.Snapshot()

	expected := uint64(workers * samplesPerWorker)

	if snapshot.LatencySamples != expected {
		t.Fatalf(
			"expected %d latency samples, got %d",
			expected,
			snapshot.LatencySamples,
		)
	}

	if snapshot.TotalLatencyMs != expected {
		t.Fatalf(
			"expected %d total latency ms, got %d",
			expected,
			snapshot.TotalLatencyMs,
		)
	}
}

func TestPerformanceMetricsAverageLatencyEmpty(t *testing.T) {
	registry := NewRegistry()

	snapshot := registry.Snapshot()

	if got := snapshot.AverageLatency(); got != 0 {
		t.Fatalf("expected zero average latency, got %v", got)
	}
}

func TestPerformanceMetricsTimestamp(t *testing.T) {
	before := time.Now()

	registry := NewRegistry()
	registry.RecordLatency(1)

	snapshot := registry.Snapshot()

	after := time.Now()

	if snapshot.UpdatedAt.Before(before) {
		t.Fatalf("metrics timestamp predates operation")
	}

	if snapshot.UpdatedAt.After(after) {
		t.Fatalf("metrics timestamp is in the future")
	}
}
