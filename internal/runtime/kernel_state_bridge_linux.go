// SPDX-License-Identifier: Apache-2.0
//go:build linux

package runtime

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

// KernelStateBridge connects the real Linux rtnetlink event stream
// to the production state enforcer.
//
// Data flow:
//
//	rtnetlink Event
//	     ↓
//	FaultEvent
//	     ↓
//	ProductionKernelEnforcer
//
// The bridge intentionally does not invent TCP state from raw
// rtnetlink messages. TCP-specific events must come from a real
// TCP event source/parser.
type KernelStateBridge struct {
	logger   *slog.Logger
	listener *kernel.NetlinkListener
	enforcer *kernel.ProductionKernelEnforcer
}

// NewKernelStateBridge creates the Linux kernel state bridge.
func NewKernelStateBridge(
	logger *slog.Logger,
	listener *kernel.NetlinkListener,
	enforcer *kernel.ProductionKernelEnforcer,
) (*KernelStateBridge, error) {
	if listener == nil {
		return nil, fmt.Errorf("kernel state bridge: listener is required")
	}

	if enforcer == nil {
		return nil, fmt.Errorf("kernel state bridge: enforcer is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &KernelStateBridge{
		logger:   logger,
		listener: listener,
		enforcer: enforcer,
	}, nil
}

// Start starts consuming real kernel events.
//
// The listener itself is started separately so lifecycle ownership
// remains explicit in the daemon.
func (b *KernelStateBridge) Start(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("kernel state bridge: context is required")
	}

	b.logger.Info("kernel state bridge started")

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("kernel state bridge stopped")
			return nil

		case event := <-b.listener.Events():
			faultEvent := faultEventFromNetlink(event)

			if err := b.enforcer.Handle(ctx, faultEvent); err != nil {
				b.logger.Error(
					"kernel state enforcement failed",
					"event_type", event.Type,
					"error", err,
				)
			}
		}
	}
}

// faultEventFromNetlink converts a real rtnetlink observation into
// the internal FaultEvent representation.
//
// This is intentionally an observation event. It does not claim that
// a topology event is a TCP connection failure.
func faultEventFromNetlink(event kernel.Event) kernel.FaultEvent {
	metadata := map[string]string{
		"source":     "rtnetlink",
		"event_type": event.Type,
	}

	if len(event.Data) > 0 {
		metadata["payload_bytes"] = fmt.Sprintf("%d", len(event.Data))
	}

	return kernel.FaultEvent{
		Type:     "kernel_warning",
		Source:   "kernel/rtnetlink",
		Reason:   event.Type,
		Metadata: metadata,
	}
}
