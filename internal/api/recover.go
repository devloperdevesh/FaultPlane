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
	){


		var req recoverRequest


		json.NewDecoder(
			r.Body,
		).Decode(&req)



		err:=controller.Recover(
			req.WorkflowID,
		)


		if err!=nil{

			http.Error(
				w,
				err.Error(),
				500,
			)

			return
		}



		json.NewEncoder(w).Encode(
			map[string]string{
				"status":"recovered",
			},
		)

	}

}