package runtime

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/config"
	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/gateway"
	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

type Daemon struct {
	logger    *slog.Logger
	control   *control.Manager
	gateway   *gateway.Manager
	kernel    *kernel.Monitor
	bpfLoader *kernel.Loader
}

func New(
	logger *slog.Logger,
	controlManager *control.Manager,
	gatewayManager *gateway.Manager,
) *Daemon {
	return &Daemon{
		logger:    logger,
		control:   controlManager,
		gateway:   gatewayManager,
		kernel:    kernel.NewMonitor(logger),
		bpfLoader: kernel.NewLoader(),
	}
}

func (d *Daemon) Start(ctx context.Context) error {
	d.logger.Info("faultplane daemon starting")

	// Existing kernel monitor remains the gateway/API event source.
	if err := d.kernel.Start(ctx); err != nil {
		return err
	}

	runtimeConfig := config.Load()

	if runtimeConfig.BPFObjectPath == "" {
		d.kernel.Stop()
		return fmt.Errorf("BPF object path is empty")
	}

	if err := d.bpfLoader.Load(runtimeConfig.BPFObjectPath); err != nil {
		d.kernel.Stop()
		return fmt.Errorf("load production eBPF programs: %w", err)
	}
	// Platform-specific kernel state enforcement.
	kernelStateRuntime := newKernelStateRuntime(d.logger, d.bpfLoader)

	if err := kernelStateRuntime.Start(ctx); err != nil {
		_ = d.bpfLoader.Close()
		d.kernel.Stop()
		return fmt.Errorf("start kernel state runtime: %w", err)
	}

	workerRegistry := NewWorkerRegistry(
		d.gateway.WorkerStore(),
		d.gateway.Registry(),
		d.gateway.Topology(),
	)

	go workerRegistry.Start(ctx)

	go func() {
		if err := d.control.Start(ctx); err != nil {
			d.logger.Error(
				"control manager failed",
				"error", err,
			)
		}
	}()

	go func() {
		if err := d.gateway.StartWithKernel(ctx, d.kernel); err != nil {
			d.logger.Error(
				"gateway manager failed",
				"error", err,
			)
		}
	}()

	<-ctx.Done()

	// --------------------------------------------------------
	// Graceful shutdown in reverse dependency order.
	// --------------------------------------------------------
	kernelStateRuntime.Stop()
	if err := d.bpfLoader.Close(); err != nil {
		d.logger.Error(
			"failed to close production eBPF loader",
			"error", err,
		)
	}

	d.kernel.Stop()
	transitions, errors := kernelStateRuntime.Stats()

	d.logger.Info(
		"faultplane daemon stopped",
		"kernel_enforcer_state", kernelStateRuntime.State(),
		"kernel_enforcer_transitions", transitions,
		"kernel_enforcer_errors", errors,
	)

	return nil
}
