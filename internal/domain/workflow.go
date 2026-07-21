package domain

import "time"

type WorkflowStatus string

const (
	StatusPending   WorkflowStatus = "pending"
	StatusRunning   WorkflowStatus = "running"
	StatusPaused    WorkflowStatus = "paused"
	StatusCompleted WorkflowStatus = "completed"
	StatusFailed    WorkflowStatus = "failed"
)

type Checkpoint struct {
	ID string

	WorkflowID string

	Step uint64

	Payload []byte

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
