package storage

import "errors"

var (
	// ErrWorkflowNotFound is returned when a workflow does not exist in storage.
	ErrWorkflowNotFound = errors.New("workflow not found")

	// ErrInvalidWorkflow is returned when an invalid workflow payload is provided.
	ErrInvalidWorkflow = errors.New("invalid workflow")
)
