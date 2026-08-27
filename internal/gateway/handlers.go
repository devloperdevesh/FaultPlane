package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricsResponse struct {
	Requests uint64  `json:"requests"`
	Latency  float64 `json:"latency"`
	CPU      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(
		w,
		http.StatusOK,
		HealthResponse{
			Status:    "healthy",
			Service:   "faultplane-gateway",
			Timestamp: time.Now(),
		},
	)
}

func metricsHandler(registry *telemetry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := registry.Snapshot()

		writeJSON(
			w,
			http.StatusOK,
			MetricsResponse{
				Requests: metrics.Requests,
				Latency:  metrics.AverageLatency(),
				CPU:      metrics.CPU,
				Memory:   metrics.Memory,
			},
		)
	})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
