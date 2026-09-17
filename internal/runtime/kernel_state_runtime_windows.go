package runtime

import "context"

type windowsKernelStateRuntime struct{}

func newKernelStateRuntime(_ interface{}, _ interface{}, _ interface{}) KernelStateRuntime {
	return &windowsKernelStateRuntime{}
}

func (r *windowsKernelStateRuntime) Start(ctx context.Context) error {
	if ctx == nil {
		return context.Canceled
	}
	return ctx.Err()
}

func (r *windowsKernelStateRuntime) Stop() {}

func (r *windowsKernelStateRuntime) Recover(ctx context.Context, workflowID string) error {
	if ctx == nil {
		return context.Canceled
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if workflowID == "" {
		return nil
	}

	return nil
}

func (r *windowsKernelStateRuntime) State() string {
	return "unknown"
}

func (r *windowsKernelStateRuntime) Stats() (uint64, uint64) {
	return 0, 0
}
