package storage

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrStateNotFound = errors.New("agent state not found")
)

type AgentState struct {
	AgentID      string    `json:"agent_id"`
	CurrentStep  int       `json:"current_step"`
	Context      string    `json:"context"`
	LastModified time.Time `json:"last_modified"`
}

type MemoryStore struct {
	mu     sync.RWMutex
	states map[string]AgentState
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		states: make(map[string]AgentState),
	}
}

func (m *MemoryStore) Save(state AgentState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	state.LastModified = time.Now().UTC()
	m.states[state.AgentID] = state
}

func (m *MemoryStore) Get(agentID string) (AgentState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, ok := m.states[agentID]
	if !ok {
		return AgentState{}, ErrStateNotFound
	}

	return state, nil
}

func (m *MemoryStore) Delete(agentID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.states, agentID)
}

func (m *MemoryStore) Exists(agentID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.states[agentID]
	return ok
}

func (m *MemoryStore) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.states)
}
