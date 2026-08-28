package time_travel

import (
	"errors"
	"testing"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/domain"
)

func TestFromCheckpoint(t *testing.T) {
	created := time.Unix(123, 0).UTC()

	checkpoint := &domain.Checkpoint{
		ID:         "checkpoint-42",
		WorkflowID: "workflow-1",
		Step:       42,
		Payload:    []byte("historical-state"),
		CreatedAt:  created,
	}

	snapshot, err := FromCheckpoint(checkpoint)
	if err != nil {
		t.Fatalf("from checkpoint: %v", err)
	}

	if snapshot.ID != checkpoint.ID {
		t.Fatalf("id = %q, want %q", snapshot.ID, checkpoint.ID)
	}

	if snapshot.WorkflowID != checkpoint.WorkflowID {
		t.Fatalf("workflow id = %q, want %q", snapshot.WorkflowID, checkpoint.WorkflowID)
	}

	if snapshot.Step != checkpoint.Step {
		t.Fatalf("step = %d, want %d", snapshot.Step, checkpoint.Step)
	}

	if string(snapshot.Payload) != "historical-state" {
		t.Fatalf("payload = %q, want historical-state", snapshot.Payload)
	}

	if !snapshot.CreatedAt.Equal(created) {
		t.Fatalf("created_at = %v, want %v", snapshot.CreatedAt, created)
	}
}

func TestFromCheckpointCopiesPayload(t *testing.T) {
	payload := []byte("immutable-source")

	checkpoint := &domain.Checkpoint{
		ID:         "checkpoint-1",
		WorkflowID: "workflow-1",
		Payload:    payload,
	}

	snapshot, err := FromCheckpoint(checkpoint)
	if err != nil {
		t.Fatalf("from checkpoint: %v", err)
	}

	snapshot.Payload[0] = 'X'

	if string(checkpoint.Payload) != "immutable-source" {
		t.Fatal("snapshot payload aliases checkpoint payload")
	}
}

func TestFromCheckpointRejectsNil(t *testing.T) {
	_, err := FromCheckpoint(nil)

	if !errors.Is(err, ErrNilCheckpoint) {
		t.Fatalf("expected ErrNilCheckpoint, got %v", err)
	}
}

func TestFromCheckpointRejectsInvalidWorkflow(t *testing.T) {
	_, err := FromCheckpoint(&domain.Checkpoint{
		ID:   "checkpoint-1",
		Step: 1,
	})

	if !errors.Is(err, ErrWorkflowMismatch) {
		t.Fatalf("expected ErrWorkflowMismatch, got %v", err)
	}

	_, err = FromCheckpoint(&domain.Checkpoint{
		WorkflowID: "workflow-1",
		Step:       1,
	})

	if !errors.Is(err, ErrWorkflowMismatch) {
		t.Fatalf("expected ErrWorkflowMismatch, got %v", err)
	}
}
