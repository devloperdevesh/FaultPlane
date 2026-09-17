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
	logger       *slog.Logger
	control      *control.Manager
	gateway      *gateway.Manager
	kernel       *kernel.Monitor
	kernelSource KernelEventSource
	bpfLoader    *kernel.Loader
}

func New(
	logger *slog.Logger,
	controlManager *control.Manager,
	gatewayManager *gateway.Manager,
) *Daemon {
	return &Daemon{
		logger:       logger,
		control:      controlManager,
		gateway:      gatewayManager,
		kernel:       kernel.NewMonitor(logger),
		kernelSource: NewKernelEventSource(logger),
		bpfLoader:    kernel.NewLoader(),
	}
}

func (d *Daemon) Start(ctx context.Context) error {
	d.logger.Info("faultplane daemon starting")

	if err := d.kernel.Start(ctx); err != nil {
		return err
	}

	// Linux: starts the real rtnetlink source.
	// Non-Linux: starts the platform-safe no-op source.
	if err := d.kernelSource.Start(ctx); err != nil {
		d.kernel.Stop()
		return fmt.Errorf("start kernel event source: %w", err)
	}

	go d.consumeKernelEvents(ctx)

	runtimeConfig := config.Load()

	if runtimeConfig.BPFObjectPath == "" {
		d.kernelSource.Stop()
		d.kernel.Stop()
		return fmt.Errorf("BPF object path is empty")
	}

	if err := d.bpfLoader.Load(runtimeConfig.BPFObjectPath); err != nil {
		d.kernelSource.Stop()
		d.kernel.Stop()
		return fmt.Errorf("load production eBPF programs: %w", err)
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

	if err := d.bpfLoader.Close(); err != nil {
		d.logger.Error(
			"failed to close production eBPF loader",
			"error", err,
		)
	}

	d.kernelSource.Stop()
	d.kernel.Stop()

	d.logger.Info("faultplane daemon stopped")

	return nil
}

// consumeKernelEvents keeps the platform kernel event source alive.
//
// On Linux this consumes real rtnetlink events.
// On non-Linux platforms the event source is intentionally a no-op.
func (d *Daemon) consumeKernelEvents(ctx context.Context) {
	events := d.kernelSource.Events()

	if events == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-events:
			if !ok {
				return
			}

			d.logger.Info(
				"kernel event entered runtime",
				"type", event.Type,
				"source", "kernel-event-source",
			)
		}
	}
}
