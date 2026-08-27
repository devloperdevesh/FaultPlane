package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/config"
	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/gateway"
	"github.com/devloperdevesh/FaultPlane/internal/logging"
	"github.com/devloperdevesh/FaultPlane/internal/runtime"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func main() {
	cfg := config.Load()

	logger := logging.New(cfg.LogLevel)

	logger.Info("starting FaultPlane daemon")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	registry := telemetry.NewRegistry()
	cpuSampler := runtime.NewProcessCPUSampler()
	registry.SetCPUSampler(cpuSampler)
	collector := telemetry.NewCollector(registry)

	controlManager := control.New(logger)

	gatewayManager := gateway.New(
		logger,
		registry,
		collector,
	)

	daemon := runtime.New(
		logger,
		controlManager,
		gatewayManager,
	)

	go func() {
		if err := daemon.Start(ctx); err != nil {
			logger.Error(
				"daemon failed",
				"error", err,
			)
			cancel()
		}
	}()

	<-ctx.Done()

	logger.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	_ = shutdownCtx

	logger.Info("FaultPlane daemon stopped cleanly")
}
