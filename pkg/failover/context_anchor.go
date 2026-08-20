// SPDX-License-Identifier: Apache-2.0
package failover

import (
	"sync/atomic"
	"unsafe"
)

// ContextAnchor provides an atomic pointer to the currently active
// connection context and tracks connection-level failover events.
type ContextAnchor struct {
	active unsafe.Pointer
	faults atomic.Uint64
}

// NewContextAnchor creates an empty context anchor.
func NewContextAnchor() *ContextAnchor {
	return &ContextAnchor{}
}

// Load returns the currently active context pointer.
func (a *ContextAnchor) Load() unsafe.Pointer {
	return atomic.LoadPointer(&a.active)
}

// Swap atomically replaces the active context.
//
// The operation succeeds only when the current pointer matches oldContext.
func (a *ContextAnchor) Swap(
	oldContext unsafe.Pointer,
	newContext unsafe.Pointer,
) bool {
	if !atomic.CompareAndSwapPointer(&a.active, oldContext, newContext) {
		return false
	}

	a.faults.Add(1)
	return true
}

// FaultCount returns the number of successful context transitions.
func (a *ContextAnchor) FaultCount() uint64 {
	return a.faults.Load()
}
