package storage

import (
	"context"
	"sync"

	"github.com/devloperdevesh/agentmesh/internal/domain"
)

// MemoryStore provides an in-memory workflow persistence backend.
// It is thread-safe and intended for local development/testing
// and lightweight runtime deployments.
type MemoryStore struct {
	mu sync.RWMutex

	workflows map[string]*domain.Workflow
}

// NewMemoryStore creates a new memory-backed storage engine.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		workflows: make(map[string]*domain.Workflow),
	}
}

// Save stores or updates a workflow.
func (m *MemoryStore) Save(
	ctx context.Context,
	workflow *domain.Workflow,
) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	if workflow == nil {
		return ErrInvalidWorkflow
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.workflows[workflow.ID] = workflow

	return nil
}

// Load retrieves a workflow by ID.
func (m *MemoryStore) Load(
	ctx context.Context,
	id string,
) (*domain.Workflow, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	workflow, ok := m.workflows[id]

	if !ok {
		return nil, ErrWorkflowNotFound
	}

	return workflow, nil
}

// Delete removes a workflow from storage.
func (m *MemoryStore) Delete(
	ctx context.Context,
	id string,
) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.workflows, id)

	return nil
}

// List returns all stored workflows.
func (m *MemoryStore) List(
	ctx context.Context,
) ([]*domain.Workflow, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	workflows := make([]*domain.Workflow, 0, len(m.workflows))

	for _, workflow := range m.workflows {
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}
