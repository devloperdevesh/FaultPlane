// SPDX-License-Identifier: Apache-2.0

package sync

import (
	"sync/atomic"
	"unsafe"
)

// HardenedMemoryFence provides explicit atomic memory-ordering operations.
//
// The implementation relies on Go's sync/atomic package for the required
// happens-before guarantees. It does not expose architecture-specific
// CPU instructions directly.
type HardenedMemoryFence struct {
	state atomic.Uint64
}

// Store publishes a synchronized state.
func (f *HardenedMemoryFence) Store(state uint64) {
	f.state.Store(state)
}

// Load returns the currently published state.
func (f *HardenedMemoryFence) Load() uint64 {
	return f.state.Load()
}

// CompareAndSwap atomically changes the state from oldState to newState.
func (f *HardenedMemoryFence) CompareAndSwap(
	oldState, newState uint64,
) bool {
	return f.state.CompareAndSwap(oldState, newState)
}

// LockFreePointerCAS atomically replaces a pointer when its current value
// matches oldMapping.
func LockFreePointerCAS(
	targetPointer *unsafe.Pointer,
	oldMapping, newMapping unsafe.Pointer,
) bool {
	return atomic.CompareAndSwapPointer(
		targetPointer,
		oldMapping,
		newMapping,
	)
}
