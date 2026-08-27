package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type RuntimeMetricsResponse struct {
	Requests    uint64  `json:"requests"`
	Workers     uint64  `json:"workers"`
	Recoveries  uint64  `json:"recoveries"`
	Checkpoints uint64  `json:"checkpoints"`
	Latency     float64 `json:"latency"`
	CPU         float64 `json:"cpu"`
	Memory      float64 `json:"memory"`
	UpdatedAt   string  `json:"updated_at"`
}

func MetricsHandler(registry *telemetry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		metrics := registry.Snapshot()

		response := RuntimeMetricsResponse{
			Requests:    metrics.Requests,
			Workers:     metrics.Workers,
			Recoveries:  metrics.Recoveries,
			Checkpoints: metrics.Checkpoints,
			Latency:     metrics.AverageLatency(),
			CPU:         metrics.CPU,
			Memory:      metrics.Memory,
			UpdatedAt:   metrics.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	})
}
