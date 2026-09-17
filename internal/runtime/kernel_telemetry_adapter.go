//go:build linux

package runtime

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

// kernelTelemetryAdapter connects the production kernel state enforcer
// to FaultPlane's shared telemetry registry.
type kernelTelemetryAdapter struct {
	registry *telemetry.Registry

	mu sync.RWMutex

	lastState      string
	lastEvent      string
	lastReason     string
	lastTransition time.Time
	lastError      time.Time
}

// NewKernelTelemetryAdapter creates a production kernel telemetry adapter.
func NewKernelTelemetryAdapter(registry *telemetry.Registry) *kernelTelemetryAdapter {
	return &kernelTelemetryAdapter{
		registry: registry,
	}
}

// RecordTransition implements kernel.StateTelemetry.
func (a *kernelTelemetryAdapter) RecordTransition(
	transition kernel.StateTransition,
) {
	if a == nil || a.registry == nil {
		return
	}

	a.mu.Lock()
	a.lastState = transition.To.String()
	a.lastEvent = transition.EventType
	a.lastReason = transition.Reason
	a.lastTransition = transition.At
	a.mu.Unlock()

	a.registry.RecordKernelTransition(
		transition.From.String(),
		transition.To.String(),
		string(transition.Action),
		transition.EventType,
		transition.Reason,
	)
}

// RecordError implements kernel.StateTelemetry.
func (a *kernelTelemetryAdapter) RecordError(err error) {
	if a == nil || a.registry == nil || err == nil {
		return
	}

	a.mu.Lock()
	a.lastError = time.Now()
	a.mu.Unlock()

	a.registry.RecordKernelError(err)
}

// LastState returns the latest committed kernel state.
func (a *kernelTelemetryAdapter) LastState() string {
	if a == nil {
		return ""
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastState
}

// LastEvent returns the latest kernel event type.
func (a *kernelTelemetryAdapter) LastEvent() string {
	if a == nil {
		return ""
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastEvent
}

// LastReason returns the latest kernel transition reason.
func (a *kernelTelemetryAdapter) LastReason() string {
	if a == nil {
		return ""
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastReason
}

// LastTransitionTime returns the latest transition timestamp.
func (a *kernelTelemetryAdapter) LastTransitionTime() time.Time {
	if a == nil {
		return time.Time{}
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastTransition
}

// LastErrorTime returns the latest kernel error timestamp.
func (a *kernelTelemetryAdapter) LastErrorTime() time.Time {
	if a == nil {
		return time.Time{}
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.lastError
}

// Validate verifies that the adapter has a production telemetry registry.
func (a *kernelTelemetryAdapter) Validate() error {
	if a == nil {
		return errors.New("kernel telemetry adapter is nil")
	}

	if a.registry == nil {
		return fmt.Errorf("kernel telemetry adapter: registry is required")
	}

	return nil
}

// Compile-time contract verification.
var _ kernel.StateTelemetry = (*kernelTelemetryAdapter)(nil)
