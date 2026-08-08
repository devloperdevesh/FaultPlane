package gateway

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`

	Service string `json:"service"`

	Timestamp time.Time `json:"timestamp"`
}

type MetricsResponse struct {
	Requests uint64 `json:"requests"`

	Checkpoints uint64 `json:"checkpoints"`

	Recoveries uint64 `json:"recoveries"`

	LatencyMs uint64 `json:"latency_ms"`
}

type WorkersResponse struct {
	Workers int `json:"workers"`

	Status string `json:"status"`
}

func healthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	writeJSON(
		w,
		http.StatusOK,
		HealthResponse{

			Status: "healthy",

			Service: "agentmesh-gateway",

			Timestamp: time.Now(),
		},
	)
}

func metricsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	writeJSON(
		w,
		http.StatusOK,
		MetricsResponse{},
	)
}

func workersHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	writeJSON(
		w,
		http.StatusOK,
		WorkersResponse{

			Workers: 0,

			Status: "idle",
		},
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value interface{},
) {

	w.Header().
		Set(
			"Content-Type",
			"application/json",
		)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).
		Encode(value)
}
