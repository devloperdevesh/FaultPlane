package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/agentmesh/internal/control"
)

type recoverRequest struct {
	WorkflowID string `json:"workflow_id"`
}

func RecoverHandler(
	controller *control.Controller,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		var req recoverRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

			http.Error(
				w,
				"invalid recovery request",
				http.StatusBadRequest,
			)

			return
		}

		if err := controller.Recover(
			req.WorkflowID,
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
				"status": "recovered",
			},
		); err != nil {

			http.Error(
				w,
				"failed to encode recovery response",
				http.StatusInternalServerError,
			)

			return
		}

	}

}
