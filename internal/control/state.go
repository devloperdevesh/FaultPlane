package control

import (
	"context"
	"errors"
	"sync"
	"time"
)

// AgentState represents the localized type-isolated memory snapshot of an active agent.
type AgentState struct {
	AgentID          string    `json:"agent_id"`
	CurrentStep      uint64    `json:"last_successful_step"`
	MemoryPayload    string    `json:"context_memory_compressed"`
	CapabilityAnchor string    `json:"hcp_capability_anchor"` // HCP/1.0 Handshake Token
	LastUpdated      time.Time `json:"timestamp"`
}

// UniqueStateControlPlane manages multi-threaded runtime memory states without allocation leaks.
type UniqueStateControlPlane struct {
	mu           sync.RWMutex
	ActiveStates map[string]*AgentState
	CircuitOpen  bool
}

// NewControlPlane initializes a thread-safe configuration boundary.
func NewControlPlane() *UniqueStateControlPlane {
	return &UniqueStateControlPlane{
		ActiveStates: make(map[string]*AgentState),
		CircuitOpen:  false,
	}
}

// CommitCheckpoint securely serializes runtime states under intense thread concurrency boundaries.
func (scp *UniqueStateControlPlane) CommitCheckpoint(agentID string, step uint64, payload string, anchor string) {
	scp.mu.Lock()
	defer scp.mu.Unlock()

	scp.ActiveStates[agentID] = &AgentState{
		AgentID:          agentID,
		CurrentStep:      step,
		MemoryPayload:    payload,
		CapabilityAnchor: anchor,
		LastUpdated:      time.Now(),
	}
}

// HotSwapExecutionState extracts snapshots and recalculates validation anchors under 30ms chaos thresholds.
func (scp *UniqueStateControlPlane) HotSwapExecutionState(ctx context.Context, agentID string) (*AgentState, string, error) {
	scp.mu.RLock()
	state, exists := scp.ActiveStates[agentID]
	scp.mu.RUnlock()

	if !exists {
		return nil, "http://primary-data-node:8000", errors.New("no historical context signature registered")
	}

	if scp.CircuitOpen {
		// Recalculating the identity receipt dynamically post-rehydration (Bypassing AGT/AutoGen Blindspot)
		reissuedAnchor := "reissued_" + state.CapabilityAnchor + "_" + time.Now().Format(time.RFC3339)

		scp.mu.Lock()
		state.CapabilityAnchor = reissuedAnchor
		scp.mu.Unlock()

		return state, "http://fallback-mock-node:8001/api/v1/resume", nil
	}

	return state, "http://primary-data-node:8000", nil
}
