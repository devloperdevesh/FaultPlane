//go:build windows

package runtime

import (
	"context"
)

type windowsKernelStateRuntime struct{}

func newKernelStateRuntime(_ interface{}, _ interface{}) KernelStateRuntime {
	return &windowsKernelStateRuntime{}
}

func (r *windowsKernelStateRuntime) Start(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}

func (r *windowsKernelStateRuntime) Stop() {}

func (r *windowsKernelStateRuntime) State() string {
	return "unsupported"
}

func (r *windowsKernelStateRuntime) Stats() (uint64, uint64) {
	return 0, 0
}
