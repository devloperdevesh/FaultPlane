package telemetry

import (
	"context"
	"testing"
	"time"
)

func TestPipelineProcessesEvents(t *testing.T) {
	registry := NewRegistry()
	pipeline := NewPipeline(registry, 8)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline.Start(ctx)

	if err := pipeline.Publish(Event{
		Name: "request",
	}); err != nil {
		t.Fatalf("publish request: %v", err)
	}

	if err := pipeline.Publish(Event{
		Name:  "latency",
		Value: 25,
	}); err != nil {
		t.Fatalf("publish latency: %v", err)
	}

	// Give the worker a chance to consume the buffered events.
	time.Sleep(20 * time.Millisecond)

	metrics := registry.Snapshot()

	if metrics.Requests != 1 {
		t.Fatalf("expected 1 request, got %d", metrics.Requests)
	}

	if metrics.LatencySamples != 1 {
		t.Fatalf(
			"expected 1 latency sample, got %d",
			metrics.LatencySamples,
		)
	}

	if metrics.TotalLatencyMs != 25 {
		t.Fatalf(
			"expected 25ms total latency, got %d",
			metrics.TotalLatencyMs,
		)
	}
}

func TestPipelineRejectsPublishBeforeStart(t *testing.T) {
	pipeline := NewPipeline(NewRegistry(), 4)

	if err := pipeline.Publish(Event{Name: "request"}); err == nil {
		t.Fatal("expected publish to fail before pipeline starts")
	}
}

func TestPipelineStopsWithContext(t *testing.T) {
	pipeline := NewPipeline(NewRegistry(), 4)

	ctx, cancel := context.WithCancel(context.Background())
	pipeline.Start(ctx)

	cancel()
	pipeline.Wait()
}
