package failover

// ConnectionState represents lifecycle state
// of a FaultPlane connection/session.
type ConnectionState string

const (
	// StatePrimary represents the initial active route.
	StatePrimary ConnectionState = "primary"

	// StateHealthy represents normal operation.
	StateHealthy ConnectionState = "healthy"

	// StateDegraded represents partial failure.
	StateDegraded ConnectionState = "degraded"

	// StateFailed represents unavailable connection.
	StateFailed ConnectionState = "failed"

	// StateRecovering represents recovery in progress.
	StateRecovering ConnectionState = "recovering"

	// StateFallback represents traffic moved to backup route.
	StateFallback ConnectionState = "fallback"
)

// Valid checks whether a state is supported.
func (s ConnectionState) Valid() bool {

	switch s {

	case StatePrimary,
		StateHealthy,
		StateDegraded,
		StateFailed,
		StateRecovering,
		StateFallback:

		return true
	}

	return false
}
