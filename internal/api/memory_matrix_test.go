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
}
