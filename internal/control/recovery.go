package control

import (
	"fmt"
	"time"
)

// Recover restores a workflow from its latest checkpoint
// and resumes execution from the saved state.
func (c *Controller) Recover(
	id string,
) error {

	start := time.Now()

	// Load latest checkpoint
	checkpoint, err := c.RestoreCheckpoint(id)

	if err != nil {

		return fmt.Errorf(
			"failed to restore checkpoint for workflow %s: %w",
			id,
			err,
		)
	}

	// Resume workflow execution
	err = c.Resume(
		checkpoint.WorkflowID,
	)

	if err != nil {

		return fmt.Errorf(
			"failed to resume workflow %s: %w",
			checkpoint.WorkflowID,
			err,
		)
	}

	// Record recovery metrics
	if c.telemetry != nil {

		c.telemetry.RecordRecovery(
			time.Since(start),
		)
	}

	return nil
}
