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
	){


		var req checkpointRequest


		err:=json.NewDecoder(
			r.Body,
		).Decode(&req)


		if err!=nil{

			http.Error(
				w,
				err.Error(),
				400,
			)

			return
		}


		err=controller.CreateCheckpoint(
			req.WorkflowID,
			req.Step,
			req.Payload,
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
				"status":"checkpoint_saved",
			},
		)

	}

}