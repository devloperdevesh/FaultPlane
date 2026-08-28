package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/control"
)

type workflowCreateRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type workflowResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	CurrentStep uint64 `json:"current_step"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func workflowResponseFrom(wf *control.Workflow) workflowResponse {
	return workflowResponse{
		ID:          wf.ID,
		Name:        wf.Name,
		Status:      string(wf.Status),
		CurrentStep: wf.CurrentStep,
		CreatedAt:   wf.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:   wf.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}

func WorkflowHandler(
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

		var req workflowCreateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(
				w,
				"invalid workflow request",
				http.StatusBadRequest,
			)
			return
		}

		if req.ID == "" {
			http.Error(
				w,
				"id is required",
				http.StatusBadRequest,
			)
			return
		}

		if req.Name == "" {
			req.Name = req.ID
		}

		workflow, err := controller.Start(
			req.ID,
			req.Name,
		)
		if err != nil {
			status := http.StatusInternalServerError

			if err == control.ErrWorkflowExists {
				status = http.StatusConflict
			}

			http.Error(
				w,
				err.Error(),
				status,
			)
			return
		}

		writeJSON(
			w,
			http.StatusCreated,
			workflowResponseFrom(workflow),
		)
	}
}
