package control

import "errors"

var (
	ErrWorkflowNotFound =
		errors.New("workflow not found")

	ErrWorkflowExists =
		errors.New("workflow already exists")

	ErrCheckpointNotFound =
		errors.New("checkpoint not found")
)