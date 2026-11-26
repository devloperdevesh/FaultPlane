package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devloperdevesh/FaultPlane/internal/storage"
)

func TestMemoryMatrixHandler(t *testing.T) {
	pool := storage.NewBufferPool(32)
	_, _ = pool.AcquireSlot("test-tenant", 1)

	handler := MemoryMatrixHandler(pool)
	req := httptest.NewRequest(http.MethodGet, "/api/memory/matrix", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var raw map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to decode JSON map: %v", err)
	}
	if val, ok := raw["l1_l2_thrashing_rate"]; !ok || val != nil {
		t.Errorf("expected l1_l2_thrashing_rate to be JSON null when unavailable, got %v (present=%v)", val, ok)
	}

	var snap storage.BufferPoolSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if snap.TotalSlots != 32 {
		t.Errorf("expected 32 slots, got %d", snap.TotalSlots)
	}

	if snap.CacheLineSize != 64 {
		t.Errorf("expected cache line size 64, got %d", snap.CacheLineSize)
	}

	if !snap.CacheAlignmentCongruent {
		t.Errorf("expected cache alignment to be congruent")
	}

	if snap.L1L2ThrashingAvailable || snap.L1L2ThrashingRate != nil {
		t.Errorf("expected L1/L2 thrashing telemetry to report unavailable")
	}

	if !snap.GCPausesAvailable {
		t.Errorf("expected GC pause telemetry to report available")
	}

	if snap.HardwareMPKEnforced || snap.IsolationMode != storage.IsolationModeLogicalTag {
		t.Errorf("expected logical software isolation status, got enforced=%v mode=%q", snap.HardwareMPKEnforced, snap.IsolationMode)
	}
}
