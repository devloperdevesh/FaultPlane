package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/domain"
)

type checkpointRequest struct {
	WorkflowID string `json:"workflow_id"`
	Step       uint64 `json:"step"`
	Payload    []byte `json:"payload"`
}

type CheckpointResponse struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflow_id"`
	Step       uint64 `json:"step"`
	Size       int    `json:"size"`
	CreatedAt  string `json:"created_at"`
}

func checkpointResponse(cp *domain.Checkpoint) CheckpointResponse {
	return CheckpointResponse{
		ID:         cp.ID,
		WorkflowID: cp.WorkflowID,
		Step:       cp.Step,
		Size:       len(cp.Payload),
		CreatedAt:  cp.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}

func CheckpointHandler(
	controller *control.Controller,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		var req checkpointRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(
				w,
				"invalid checkpoint request",
				http.StatusBadRequest,
			)
			return
		}

		if req.WorkflowID == "" {
			http.Error(
				w,
				"workflow_id is required",
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

		checkpoint, err := controller.RestoreCheckpoint(req.WorkflowID)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		writeJSON(
			w,
			http.StatusOK,
			checkpointResponse(checkpoint),
		)
	}
}
