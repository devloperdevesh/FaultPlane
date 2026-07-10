package agent_sim

import (
	"fmt"
	"time"
)


type Agent struct {

	Workflow Workflow

	Checkpoint int

}


func NewAgent(id string) *Agent {

	return &Agent{

		Workflow: Workflow{

			ID:id,

			Status:"running",

			CurrentStep:0,

			UpdatedAt:time.Now(),

		},

	}

}



func (a *Agent) ExecuteStep(){

	a.Workflow.CurrentStep++

	a.Workflow.UpdatedAt=time.Now()


	fmt.Printf(
		"[Agent] executed step %d\n",
		a.Workflow.CurrentStep,
	)

}



func (a *Agent) CreateCheckpoint(){

	a.Checkpoint=a.Workflow.CurrentStep


	fmt.Printf(
		"[Checkpoint] saved at step %d\n",
		a.Checkpoint,
	)

}