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
)

func main() {

	// Configuration
	cfg := config.Load()

	// Logger
	logger := logging.New(
		cfg.LogLevel,
	)

	logger.Info(
		"starting FaultPlane daemon",
	)

	// Shutdown context
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()

	// Control manager
	controlManager := control.New(
		logger,
	)

	// Gateway manager
	gatewayManager := gateway.New(
		logger,
	)

	// Runtime daemon
	daemon := runtime.New(
		logger,
		controlManager,
		gatewayManager,
	)

	// Start daemon
	go func() {

		if err := daemon.Start(ctx); err != nil {

			logger.Error(
				"daemon failed",
				"error",
				err,
			)

			cancel()
		}

	}()

	<-ctx.Done()

	logger.Info(
		"shutdown signal received",
	)

	shutdownCtx, shutdownCancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

	defer shutdownCancel()

	_ = shutdownCtx

	logger.Info(
		"FaultPlane daemon stopped cleanly",
	)
}
