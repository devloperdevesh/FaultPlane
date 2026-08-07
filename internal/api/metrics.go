package api

import (
	"encoding/json"
	"net/http"
)

func MetricsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	response := map[string]int{
		"active_workflows": 0,
		"recoveries":       0,
		"checkpoints":      0,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {

		http.Error(
			w,
			"failed to encode metrics response",
			http.StatusInternalServerError,
		)

		return
	}

}
