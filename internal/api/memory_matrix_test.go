package api

import (
	"bytes"
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

	if !snap.PCIeDMABound {
		t.Errorf("expected PCIe DMA to be bound")
	}

	if !snap.DPDKBypassActive {
		t.Errorf("expected DPDK bypass to be active")
	}

	if !snap.SubTwoMsCongruent {
		t.Errorf("expected sub-2ms bounds to be congruent")
	}
}

func TestMemoryDMAHandler(t *testing.T) {
	pool := storage.NewBufferPool(32)
	handler := MemoryDMAHandler(pool)

	req := httptest.NewRequest(http.MethodGet, "/api/memory/dma", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["pcie_dma_bound"] != true {
		t.Errorf("expected pcie_dma_bound to be true")
	}

	if resp["dpdk_bypass_active"] != true {
		t.Errorf("expected dpdk_bypass_active to be true")
	}
}

func TestMemoryDMAStreamTriggerHandler(t *testing.T) {
	pool := storage.NewBufferPool(32)
	_, _ = pool.AcquireSlot("test-stream", 2)
	handler := MemoryDMAStreamTriggerHandler(pool)

	payload := []byte(`{"slot_id": 0, "expert_id": 5, "tensor_bytes": 2097152}`)
	req := httptest.NewRequest(http.MethodPost, "/api/memory/dma/stream", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	var res storage.DMATransferResult
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.SubTwoMsBoundMet {
		t.Errorf("expected sub-2ms bound met, got latency: %f ms", res.LatencyMs)
	}

	if !res.KernelBypassed {
		t.Errorf("expected kernel bypassed to be true")
	}
}
