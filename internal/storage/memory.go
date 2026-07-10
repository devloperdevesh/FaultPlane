package storage

import (
	"context"
	"sync"

	"github.com/devloperdevesh/agentmesh/internal/control"
)

type MemoryStore struct {
	mu sync.RWMutex

	workflows map[string]*control.Workflow
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		workflows: make(map[string]*control.Workflow),
	}
}

func (m *MemoryStore) Save(
	ctx context.Context,
	workflow *control.Workflow,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.workflows[workflow.ID] = workflow

	return nil
}

func (m *MemoryStore) Load(
	ctx context.Context,
	id string,
) (*control.Workflow, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	workflow, ok := m.workflows[id]

	if !ok {
		return nil, ErrWorkflowNotFound
	}

	return workflow, nil
}

func (m *MemoryStore) Delete(
	ctx context.Context,
	id string,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.workflows, id)

	return nil
}

func (m *MemoryStore) List(
	ctx context.Context,
) ([]*control.Workflow, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	workflows := make([]*control.Workflow, 0, len(m.workflows))

	for _, workflow := range m.workflows {
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}