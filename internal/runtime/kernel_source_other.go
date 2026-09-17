//go:build !linux

package runtime

import (
	"context"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

// KernelEventSource is a no-op kernel event source on non-Linux systems.
type KernelEventSource interface {
	Start(context.Context) error
	Stop()
	Events() <-chan kernel.KernelEvent
}

// NonLinuxKernelEventSource preserves the runtime lifecycle without
// pretending that rtnetlink is available on this platform.
type NonLinuxKernelEventSource struct{}

func NewKernelEventSource(logger *slog.Logger) KernelEventSource {
	_ = logger

	return &NonLinuxKernelEventSource{}
}

func (s *NonLinuxKernelEventSource) Start(ctx context.Context) error {
	_ = ctx
	return nil
}

func (s *NonLinuxKernelEventSource) Stop() {}

func (s *NonLinuxKernelEventSource) Events() <-chan kernel.KernelEvent {
	return nil
}
