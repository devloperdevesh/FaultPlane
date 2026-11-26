package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/storage"
)

var globalBufferPool *storage.BufferPool

func init() {
	globalBufferPool = storage.NewBufferPool(64)

	// Seed realistic enterprise multi-tenant workloads
	tenants := []struct {
		id     string
		domain uint8
		count  int
	}{
		{"org-finance-alpha-01", 1, 14},
		{"org-defense-beta-02", 2, 12},
		{"org-health-gamma-03", 3, 10},
		{"org-sovereign-delta-04", 4, 8},
	}

	for _, t := range tenants {
		for i := 0; i < t.count; i++ {
			slot, err := globalBufferPool.AcquireSlot(t.id, t.domain)
			if err == nil && (i%7 == 0) {
				globalBufferPool.FenceSlot(slot.SlotID)
			}
		}
	}
}

// GetGlobalBufferPool returns the active buffer pool instance.
func GetGlobalBufferPool() *storage.BufferPool {
	return globalBufferPool
}

// MemoryMatrixHandler returns a point-in-time snapshot of the cache-line buffer pool.
func MemoryMatrixHandler(pool *storage.BufferPool) http.Handler {
	if pool == nil {
		pool = globalBufferPool
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := pool.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(snap)
	})
}

// MemoryMatrixStreamHandler streams real-time memory matrix telemetry via SSE.
func MemoryMatrixStreamHandler(pool *storage.BufferPool) http.Handler {
	if pool == nil {
		pool = globalBufferPool
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ticker := time.NewTicker(33 * time.Millisecond) // ~30-60 FPS stream
		defer ticker.Stop()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				snap := pool.Snapshot()
				data, err := json.Marshal(snap)
				if err != nil {
					return
				}

				if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	})
}
