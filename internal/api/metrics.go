package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/models"
)

func MetricsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	metrics := models.Metrics{
		Requests: 120,
		Latency:  35,
		CPU:      40,
		Memory:   512,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(
			w,
			"failed to encode metrics response",
			http.StatusInternalServerError,
		)
		return
	}
}
