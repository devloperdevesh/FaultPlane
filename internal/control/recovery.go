package control

import (
	"context"
	"fmt"
	"time"
)

// Recover restores a workflow from its latest checkpoint, synchronizes
// kernel/runtime state, and resumes execution from the saved state.
func (c *Controller) Recover(id string) error {
	return c.RecoverContext(context.Background(), id)
}

// RecoverContext is the context-aware recovery path used by API callers.
func (c *Controller) RecoverContext(ctx context.Context, id string) error {
	if ctx == nil {
		return fmt.Errorf("recover workflow %s: context is nil", id)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	start := time.Now()

	checkpoint, err := c.RestoreCheckpoint(id)
	if err != nil {
		return fmt.Errorf(
			"failed to restore checkpoint for workflow %s: %w",
			id,
			err,
		)
	}

	if c.runtime != nil {
		if err := c.runtime.Recover(ctx, checkpoint.WorkflowID); err != nil {
			return fmt.Errorf(
				"failed to recover runtime for workflow %s: %w",
				checkpoint.WorkflowID,
				err,
			)
		}
	}

	if err := c.Resume(checkpoint.WorkflowID); err != nil {
		return fmt.Errorf(
			"failed to resume workflow %s: %w",
			checkpoint.WorkflowID,
			err,
		)
	}

	if c.telemetry != nil {
		c.telemetry.RecordRecovery(time.Since(start))
	}

	return nil
}
