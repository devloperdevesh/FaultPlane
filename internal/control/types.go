package control

import "time"

type WorkflowStatus string

const (
	StatusRunning   WorkflowStatus = "running"
	StatusPaused    WorkflowStatus = "paused"
	StatusFailed    WorkflowStatus = "failed"
	StatusCompleted WorkflowStatus = "completed"
)


type Checkpoint struct {
	ID        string
	WorkflowID string
	Step      uint64
	Payload   []byte
	CreatedAt time.Time
}


type Workflow struct {
	ID string

	Name string

	Status WorkflowStatus

	CurrentStep uint64

	Checkpoint *Checkpoint

	CreatedAt time.Time

	UpdatedAt time.Time
}