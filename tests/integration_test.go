// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"testing"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
	"github.com/devloperdevesh/FaultPlane/pkg/failover"
)

func TestEndToEndFailoverBehavior(t *testing.T) {
	store := failover.NewStateStore(failover.StateHealthy)
	manager := failover.NewManager(store)

	registry := telemetry.NewRegistry()
	collector := telemetry.NewCollector(registry)

	// Primary starts healthy.
	if got := store.Load(); got != failover.StateHealthy {
		t.Fatalf(
			"expected primary state %q, got %q",
			failover.StateHealthy,
			got,
		)
	}

	// Initial workload.
	collector.RecordRequest()
	collector.RecordRequest()

	// Healthy -> Degraded.
	event, err := manager.HandleFailure("primary connection degraded")
	if err != nil {
		t.Fatalf("handle first failure: %v", err)
	}

	if event.Previous != failover.StateHealthy {
		t.Fatalf(
			"expected previous state %q, got %q",
			failover.StateHealthy,
			event.Previous,
		)
	}

	if event.Current != failover.StateDegraded {
		t.Fatalf(
			"expected degraded state %q, got %q",
			failover.StateDegraded,
			event.Current,
		)
	}

	if got := store.Load(); got != failover.StateDegraded {
		t.Fatalf(
			"expected store state %q, got %q",
			failover.StateDegraded,
			got,
		)
	}

	// Degraded -> Failed.
	event, err = manager.HandleFailure("primary connection lost")
	if err != nil {
		t.Fatalf("handle primary failure: %v", err)
	}

	if event.Previous != failover.StateDegraded {
		t.Fatalf(
			"expected previous state %q, got %q",
			failover.StateDegraded,
			event.Previous,
		)
	}

	if event.Current != failover.StateFailed {
		t.Fatalf(
			"expected failed state %q, got %q",
			failover.StateFailed,
			event.Current,
		)
	}

	// Failed -> Fallback.
	event, err = manager.HandleFailure("primary unavailable")
	if err != nil {
		t.Fatalf("select fallback: %v", err)
	}

	if event.Previous != failover.StateFailed {
		t.Fatalf(
			"expected previous state %q, got %q",
			failover.StateFailed,
			event.Previous,
		)
	}

	if event.Current != failover.StateFallback {
		t.Fatalf(
			"expected fallback state %q, got %q",
			failover.StateFallback,
			event.Current,
		)
	}

	if got := store.Load(); got != failover.StateFallback {
		t.Fatalf(
			"expected store state %q, got %q",
			failover.StateFallback,
			got,
		)
	}

	// Workload continues through fallback.
	collector.RecordRequest()
	collector.RecordRequest()
	collector.RecordCheckpoint()

	// Fallback -> Recovering.
	event, err = manager.Recover("primary recovery initiated")
	if err != nil {
		t.Fatalf("start recovery: %v", err)
	}

	if event.Previous != failover.StateFallback {
		t.Fatalf(
			"expected previous state %q, got %q",
			failover.StateFallback,
			event.Previous,
		)
	}

	if event.Current != failover.StateRecovering {
		t.Fatalf(
			"expected recovering state %q, got %q",
			failover.StateRecovering,
			event.Current,
		)
	}

	// Recovering -> Healthy.
	event, err = manager.Recover("primary restored")
	if err != nil {
		t.Fatalf("complete recovery: %v", err)
	}

	if event.Previous != failover.StateRecovering {
		t.Fatalf(
			"expected previous state %q, got %q",
			failover.StateRecovering,
			event.Previous,
		)
	}

	if event.Current != failover.StateHealthy {
		t.Fatalf(
			"expected final healthy state %q, got %q",
			failover.StateHealthy,
			event.Current,
		)
	}

	if got := store.Load(); got != failover.StateHealthy {
		t.Fatalf(
			"expected final state %q, got %q",
			failover.StateHealthy,
			got,
		)
	}

	// Verify workload metrics.
	metrics := registry.Snapshot()

	if metrics.Requests != 4 {
		t.Fatalf(
			"expected 4 recorded requests, got %d",
			metrics.Requests,
		)
	}

	if metrics.Checkpoints != 1 {
		t.Fatalf(
			"expected 1 checkpoint, got %d",
			metrics.Checkpoints,
		)
	}

	if metrics.Recoveries != 0 {
		t.Fatalf(
			"expected 0 recovery metrics before explicit recovery recording, got %d",
			metrics.Recoveries,
		)
	}

	// Record completed recovery latency.
	collector.RecordRecovery(25 * time.Millisecond)

	metrics = registry.Snapshot()

	if metrics.Recoveries != 1 {
		t.Fatalf(
			"expected 1 recovery, got %d",
			metrics.Recoveries,
		)
	}

	if metrics.LatencySamples != 1 {
		t.Fatalf(
			"expected 1 latency sample, got %d",
			metrics.LatencySamples,
		)
	}

	if metrics.TotalLatencyMs != 25 {
		t.Fatalf(
			"expected 25ms total recovery latency, got %dms",
			metrics.TotalLatencyMs,
		)
	}

	if metrics.AverageLatency() != 25 {
		t.Fatalf(
			"expected 25ms average latency, got %.2fms",
			metrics.AverageLatency(),
		)
	}

	t.Logf(
		"end-to-end failover verified: final_state=%s requests=%d checkpoints=%d recoveries=%d latency_ms=%d",
		store.Load(),
		metrics.Requests,
		metrics.Checkpoints,
		metrics.Recoveries,
		metrics.TotalLatencyMs,
	)
}
