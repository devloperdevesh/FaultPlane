package agent_sim


import "fmt"


func (a *Agent) Crash(){

	a.Workflow.Status="failed"


	fmt.Println(
		"[Crash] agent execution stopped",
	)

}



func (a *Agent) Recover(){

	a.Workflow.Status="running"

	a.Workflow.CurrentStep=a.Checkpoint


	fmt.Printf(
		"[Recovery] restored checkpoint step %d\n",
		a.Checkpoint,
	)

}