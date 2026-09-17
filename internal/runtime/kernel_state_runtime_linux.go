//go:build linux

package runtime

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

type linuxKernelStateRuntime struct {
	listener *kernel.NetlinkListener
	enforcer *kernel.ProductionKernelEnforcer
	bridge   *KernelStateBridge
	logger   *slog.Logger
}

func newKernelStateRuntime(logger *slog.Logger, loader *kernel.Loader) KernelStateRuntime {
	if logger == nil {
		logger = slog.Default()
	}
	if loader == nil {
		panic("create kernel state runtime: BPF loader is required")
	}

	listener := kernel.NewNetlinkListener(logger)
	telemetry := &kernel.InMemoryStateTelemetry{}

	kernelAction, err := kernel.NewBPFKernelAction(logger, loader)
	if err != nil {
		panic(fmt.Sprintf("create BPF kernel action: %v", err))
	}

	enforcer, err := kernel.NewProductionKernelEnforcer(
		logger,
		kernel.DefaultFaultPolicy{},
		kernelAction,
		telemetry,
	)
	if err != nil {
		panic(fmt.Sprintf("create production kernel enforcer: %v", err))
	}

	bridge, err := NewKernelStateBridge(logger, listener, enforcer)
	if err != nil {
		panic(fmt.Sprintf("create kernel state bridge: %v", err))
	}

	return &linuxKernelStateRuntime{
		listener: listener,
		enforcer: enforcer,
		bridge:   bridge,
		logger:   logger,
	}
}

func (r *linuxKernelStateRuntime) Start(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("start kernel state runtime: context is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := r.listener.Start(ctx); err != nil {
		return fmt.Errorf("start rtnetlink listener: %w", err)
	}

	go func() {
		if err := r.bridge.Start(ctx); err != nil && ctx.Err() == nil {
			r.logger.Error(
				"kernel state bridge stopped with error",
				"error", err,
			)
		}
	}()

	return nil
}

func (r *linuxKernelStateRuntime) Stop() {
	if r == nil {
		return
	}

	if r.listener != nil {
		r.listener.Stop()
	}
}

func (r *linuxKernelStateRuntime) State() string {
	if r == nil || r.enforcer == nil {
		return "unknown"
	}

	return r.enforcer.State().String()
}

func (r *linuxKernelStateRuntime) Stats() (uint64, uint64) {
	if r == nil || r.enforcer == nil {
		return 0, 0
	}

	return r.enforcer.Stats()
}
