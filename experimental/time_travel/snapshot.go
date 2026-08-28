package time_travel

import (
	"errors"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/domain"
)

var (
	ErrNilCheckpoint    = errors.New("time-travel checkpoint is nil")
	ErrWorkflowMismatch = errors.New("time-travel checkpoint workflow mismatch")
)

type Snapshot struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflow_id"`
	Step       uint64    `json:"step"`
	Payload    []byte    `json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
}

func FromCheckpoint(checkpoint *domain.Checkpoint) (Snapshot, error) {
	if checkpoint == nil {
		return Snapshot{}, ErrNilCheckpoint
	}

	if checkpoint.ID == "" || checkpoint.WorkflowID == "" {
		return Snapshot{}, ErrWorkflowMismatch
	}

	return Snapshot{
		ID:         checkpoint.ID,
		WorkflowID: checkpoint.WorkflowID,
		Step:       checkpoint.Step,
		Payload:    append([]byte(nil), checkpoint.Payload...),
		CreatedAt:  checkpoint.CreatedAt,
	}, nil
}
