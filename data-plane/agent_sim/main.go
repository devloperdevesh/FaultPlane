package main

import (
	"fmt"

	"github.com/devloperdevesh/agentmesh/data-plane/agent_sim/simulator"
)

func main() {
	fmt.Println("=== AgentMesh Recovery Demo ===")

	agent := simulator.NewAgent("demo-agent")

	agent.ExecuteStep()
	agent.ExecuteStep()

	agent.CreateCheckpoint()

	agent.ExecuteStep()

	agent.Crash()

	agent.Recover()

	agent.ExecuteStep()

	fmt.Println("=== Workflow Completed ===")
}
