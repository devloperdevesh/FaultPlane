package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/models"
)

func TestWorkerFailureRecoveryTelemetry(t *testing.T) {
	events := NewTelemetryStore()
	store := NewWorkerStoreWithTelemetry(events)
	topology := control.NewTopologyController()
	worker := models.Worker{ID: "worker-a", Region: "region-a", Status: "ACTIVE"}
	store.Set(worker)
	store.Set(models.Worker{ID: "worker-b", Region: "region-b", Status: "ACTIVE"})
	if len(events.Snapshot()) != 0 {
		t.Fatal("initial healthy registration is not a recovery")
	}
	worker.Status = "FAILED"
	store.Set(worker)
	store.Set(worker)
	topology.RefreshFromWorkers(store.List())
	failedResponse := httptest.NewRecorder()
	NewTopologyHandler(topology).Get(failedResponse, httptest.NewRequest("GET", "/api/topology", nil))
	var snapshot control.TopologySnapshot
	if err := json.Unmarshal(failedResponse.Body.Bytes(), &snapshot); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, node := range snapshot.Nodes {
		if node.ID == worker.ID {
			found = node.Status == "FAILED" && node.Region == "region-a"
		}
	}
	if !found {
		t.Fatal("failed regional worker absent from topology endpoint")
	}
	worker.Status = "ACTIVE"
	store.Set(worker)
	store.Set(worker)
	response := httptest.NewRecorder()
	TelemetryHandler(events).ServeHTTP(response, httptest.NewRequest("GET", "/api/telemetry", nil))
	var payload TelemetryResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Events) != 2 {
		t.Fatalf("expected one failure and recovery, got %+v", payload.Events)
	}
	for i, kind := range []string{"worker_failed", "worker_recovered"} {
		event := payload.Events[i]
		if event.Type != kind || event.Timestamp.IsZero() || event.Metadata["node_id"] != worker.ID || event.Metadata["region"] != worker.Region || event.Metadata["source"] != "worker_registry" {
			t.Fatalf("unexpected event: %+v", event)
		}
	}
	if payload.Events[1].Timestamp.Before(payload.Events[0].Timestamp) {
		t.Fatal("recovery precedes failure")
	}
	if store.List()[0].Region == "" {
		t.Fatal("worker region was lost")
	}
}

func TestWorkerHealthDoesNotInventRecovery(t *testing.T) {
	events := NewTelemetryStore()
	store := NewWorkerStoreWithTelemetry(events)
	for _, status := range []string{"FAILED", "unknown", "ACTIVE", "ACTIVE"} {
		store.Set(models.Worker{ID: "worker", Status: status})
	}
	if len(events.Snapshot()) != 1 {
		t.Fatal("unknown to healthy is not a confirmed failed-to-healthy transition")
	}
	// The legacy constructor remains usable without an observer.
	NewWorkerStore().Set(models.Worker{ID: "legacy", Status: "FAILED"})
}
