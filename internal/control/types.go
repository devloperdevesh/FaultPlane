package control

import "github.com/devloperdevesh/agentmesh/internal/domain"

// Type aliases for domain types to maintain backward compatibility
type WorkflowStatus = domain.WorkflowStatus
type Checkpoint = domain.Checkpoint
type Workflow = domain.Workflow

const (
	StatusRunning   = domain.StatusRunning
	StatusPaused    = domain.StatusPaused
	StatusFailed    = domain.StatusFailed
	StatusCompleted = domain.StatusCompleted
)
