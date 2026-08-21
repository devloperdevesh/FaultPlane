// SPDX-License-Identifier: Apache-2.0

package failover

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

var (
	ErrOverloaded      = errors.New("failover orchestration: load shedding active")
	ErrPointerMismatch = errors.New("failover orchestration: active data-plane pointer changed")
)

// Layer7OrchestrationBridge coordinates application-level load state with
// the lock-free Layer 4 failover boundary.
//
// The bridge intentionally keeps policy and transport state separate:
// Layer 7 decides whether failover work is currently permitted, while the
// existing atomic pointer performs the actual data-plane handoff.
type Layer7OrchestrationBridge struct {
	TokenRateLimit    atomic.Uint64
	AtomicLoadShedder atomic.Uint32
	DataPlaneProxy    atomic.Pointer[byte]
}

// NewInfrastructureStack creates an orchestration bridge.
//
// initialL4Proxy may be nil when the transport endpoint has not yet been
// attached.
func NewInfrastructureStack(initialL4Proxy unsafe.Pointer) *Layer7OrchestrationBridge {
	bridge := &Layer7OrchestrationBridge{}

	bridge.TokenRateLimit.Store(1_000_000)

	if initialL4Proxy != nil {
		bridge.DataPlaneProxy.Store((*byte)(initialL4Proxy))
	}

	return bridge
}

// SetLoadShedding enables or disables Layer 7 load shedding.
func (b *Layer7OrchestrationBridge) SetLoadShedding(enabled bool) {
	if enabled {
		b.AtomicLoadShedder.Store(1)
		return
	}

	b.AtomicLoadShedder.Store(0)
}

// LoadShedding reports whether failover operations are currently restricted.
func (b *Layer7OrchestrationBridge) LoadShedding() bool {
	return b.AtomicLoadShedder.Load() == 1
}

// AddTokenBudget records the amount of application-level work associated with
// the active workload.
func (b *Layer7OrchestrationBridge) AddTokenBudget(tokens uint64) {
	b.TokenRateLimit.Add(tokens)
}

// ExecuteSovereignOrchestration atomically moves the active data-plane
// reference from oldFD to fallbackFD.
//
// This function does not claim to move kernel file descriptors itself.
// It performs the lock-free state handoff used by the Go control plane.
func (b *Layer7OrchestrationBridge) ExecuteSovereignOrchestration(
	oldFD unsafe.Pointer,
	fallbackFD unsafe.Pointer,
) (bool, error) {
	if b.LoadShedding() {
		return false, ErrOverloaded
	}

	if oldFD == nil || fallbackFD == nil {
		return false, errors.New("failover orchestration: nil data-plane pointer")
	}

	swapped := b.DataPlaneProxy.CompareAndSwap(
		(*byte)(oldFD),
		(*byte)(fallbackFD),
	)
	if !swapped {
		return false, ErrPointerMismatch
	}

	b.AddTokenBudget(45_000)

	return true, nil
}

// ActiveDataPlane returns the currently published data-plane reference.
func (b *Layer7OrchestrationBridge) ActiveDataPlane() unsafe.Pointer {
	return unsafe.Pointer(b.DataPlaneProxy.Load())
}
