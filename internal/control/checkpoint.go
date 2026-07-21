package control

import (
	"context"
	"time"

	"github.com/devloperdevesh/agentmesh/internal/domain"
)

// CreateCheckpoint stores the current workflow state
// and persists it through storage backend.
func (c *Controller) CreateCheckpoint(
	id string,
	step uint64,
	payload []byte,
) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	workflow, exists := c.workflows[id]

	if !exists {
		return ErrWorkflowNotFound
	}

	checkpoint := &domain.Checkpoint{
		ID:         id,
		WorkflowID: id,
		Step:       step,
		Payload:    payload,
		CreatedAt:  time.Now(),
	}

	// Update runtime state

	workflow.Checkpoint = checkpoint

	workflow.CurrentStep = step

	workflow.UpdatedAt = time.Now()

	// Persist workflow state

	if c.storage != nil {

		err := c.storage.Save(
			context.Background(),
			workflow,
		)

		if err != nil {
			return err
		}
	}

	// Record telemetry

	if c.telemetry != nil {

		c.telemetry.RecordCheckpoint()

	}

	return nil
}

// RestoreCheckpoint retrieves latest checkpoint.
func (c *Controller) RestoreCheckpoint(
	id string,
) (*domain.Checkpoint, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	workflow, exists := c.workflows[id]

	if !exists {
		return nil, ErrWorkflowNotFound
	}

	if workflow.Checkpoint == nil {
		return nil, ErrCheckpointNotFound
	}

	return workflow.Checkpoint, nil
}
