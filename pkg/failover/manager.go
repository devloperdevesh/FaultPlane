package failover

import (
	"fmt"
)

// RecoveryEvent represents a state transition event.
type RecoveryEvent struct {
	Previous ConnectionState
	Current  ConnectionState
	Reason   string
}

// Manager controls failover decisions.
type Manager struct {
	store *StateStore
}

// NewManager creates a failover manager.
func NewManager(
	store *StateStore,
) *Manager {

	return &Manager{
		store: store,
	}
}

// HandleFailure processes detected failures.
//
// Flow:
//
// Failed detected
//
//	|
//	v
//
// Check current state
//
//	|
//	v
//
// Select fallback
//
//	|
//	v
//
// Atomic transition
func (m *Manager) HandleFailure(
	reason string,
) (RecoveryEvent, error) {

	current := m.store.Load()

	switch current {

	case StateHealthy:
		return m.transition(
			StateDegraded,
			reason,
		)

	case StateDegraded:
		return m.transition(
			StateFailed,
			reason,
		)

	case StateFailed:
		return m.transition(
			StateFallback,
			reason,
		)

	default:
		return RecoveryEvent{}, fmt.Errorf(
			"cannot failover from state %s",
			current,
		)
	}
}

// Recover moves fallback connection back.
//
// Fallback
//
//	|
//	v
//
// Recovering
//
//	|
//	v
//
// Healthy
func (m *Manager) Recover(
	reason string,
) (RecoveryEvent, error) {

	current := m.store.Load()

	switch current {

	case StateFallback:

		return m.transition(
			StateRecovering,
			reason,
		)

	case StateRecovering:

		return m.transition(
			StateHealthy,
			reason,
		)

	default:

		return RecoveryEvent{}, fmt.Errorf(
			"cannot recover from state %s",
			current,
		)
	}
}

func (m *Manager) transition(
	next ConnectionState,
	reason string,
) (RecoveryEvent, error) {

	previous := m.store.Load()

	if !next.Valid() {
		return RecoveryEvent{}, fmt.Errorf(
			"invalid target state %s",
			next,
		)
	}

	m.store.Swap(next)

	return RecoveryEvent{
		Previous: previous,
		Current:  next,
		Reason:   reason,
	}, nil
}
