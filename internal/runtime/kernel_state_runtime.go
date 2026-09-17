package runtime

import "context"

// KernelStateRuntime owns the platform-specific kernel state enforcement
// lifecycle used by the daemon and exposes explicit control-plane actions.
type KernelStateRuntime interface {
	Start(context.Context) error
	Stop()
	Recover(context.Context, string) error
	State() string
	Stats() (uint64, uint64)
}
