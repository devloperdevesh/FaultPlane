package failover

import "sync/atomic"

// StateStore provides lock-free
// connection state updates.
type StateStore struct {
	state atomic.Pointer[ConnectionState]
}

// NewStateStore creates a new state store.
func NewStateStore(
	initial ConnectionState,
) *StateStore {

	store := &StateStore{}

	store.state.Store(&initial)

	return store
}

// Load returns current connection state.
func (s *StateStore) Load() ConnectionState {

	current := s.state.Load()

	if current == nil {
		return StateFailed
	}

	return *current
}

// Swap atomically replaces current state.
func (s *StateStore) Swap(
	next ConnectionState,
) {

	s.state.Store(&next)

}
