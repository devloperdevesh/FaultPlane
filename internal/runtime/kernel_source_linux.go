//go:build linux

package runtime

import (
	"context"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

// KernelEventSource is the runtime-facing kernel event abstraction.
//
// Linux provides a real rtnetlink-backed implementation.
// Other operating systems provide a no-op implementation.
type KernelEventSource interface {
	Start(context.Context) error
	Stop()
	Events() <-chan kernel.Event
}

// LinuxKernelEventSource adapts the real Linux rtnetlink listener.
type LinuxKernelEventSource struct {
	netlink *kernel.NetlinkListener
}

// NewKernelEventSource creates the Linux kernel event source.
func NewKernelEventSource(logger *slog.Logger) KernelEventSource {
	return &LinuxKernelEventSource{
		netlink: kernel.NewNetlinkListener(logger),
	}
}

func (s *LinuxKernelEventSource) Start(ctx context.Context) error {
	return s.netlink.Start(ctx)
}

func (s *LinuxKernelEventSource) Stop() {
	if s.netlink != nil {
		s.netlink.Stop()
	}
}

func (s *LinuxKernelEventSource) Events() <-chan kernel.Event {
	if s.netlink == nil {
		return nil
	}

	return s.netlink.Events()
}
