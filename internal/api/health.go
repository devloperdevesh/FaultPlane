package api

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	response := map[string]string{
		"status":  "healthy",
		"service": "agentmesh",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {

		http.Error(
			w,
			"failed to encode health response",
			http.StatusInternalServerError,
		)

		return
	}

}
