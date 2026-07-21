package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/agentmesh/internal/control"
)

type checkpointRequest struct {
	WorkflowID string `json:"workflow_id"`

	Step uint64 `json:"step"`

	Payload []byte `json:"payload"`
}

func CheckpointHandler(
	controller *control.Controller,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		var req checkpointRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

			http.Error(
				w,
				"invalid checkpoint request",
				http.StatusBadRequest,
			)

			return
		}

		if err := controller.CreateCheckpoint(
			req.WorkflowID,
			req.Step,
			req.Payload,
		); err != nil {

			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)

			return
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		if err := json.NewEncoder(w).Encode(
			map[string]string{
				"status": "checkpoint_saved",
			},
		); err != nil {

			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)

			return
		}

	}

}
