package runtime

import "context"

// KernelStateRuntime owns the platform-specific kernel state enforcement
// lifecycle used by the daemon.
type KernelStateRuntime interface {
	Start(context.Context) error
	Stop()
	State() string
	Stats() (uint64, uint64)
}
