package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/control"
)

type recoverRequest struct {
	WorkflowID string `json:"workflow_id"`
}

func RecoverHandler(controller *control.Controller) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req recoverRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.WorkflowID == "" {
			http.Error(w, "workflow_id is required", http.StatusBadRequest)
			return
		}

		if err := controller.RecoverContext(r.Context(), req.WorkflowID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "recovered",
		})
	}
}
