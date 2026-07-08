package storage

import (
	"context"
	"sync"
)

type MemoryStore struct {
	mu           sync.RWMutex
	storePayload map[string]string
	storeSteps   map[string]uint64
	storeAnchors map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		storePayload: make(map[string]string),
		storeSteps:   make(map[string]uint64),
		storeAnchors: make(map[string]string),
	}
}

func (m *MemoryStore) Save(ctx context.Context, agentID string, step uint64, payload string, anchor string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.storePayload[agentID] = payload
	m.storeSteps[agentID] = step
	m.storeAnchors[agentID] = anchor
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, agentID string) (string, uint64, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	payload, exists := m.storePayload[agentID]
	if !exists {
		return "", 0, "", ErrStateNotFound
	}

	step := m.storeSteps[agentID]
	anchor := m.storeAnchors[agentID]
	return payload, step, anchor, nil
}
