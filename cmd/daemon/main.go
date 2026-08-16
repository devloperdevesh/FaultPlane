package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/config"
	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/gateway"
	"github.com/devloperdevesh/FaultPlane/internal/logging"
	"github.com/devloperdevesh/FaultPlane/internal/runtime"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func main() {

	// Load configuration
	cfg := config.Load()


	// Structured logger
	logger := logging.New(
		cfg.LogLevel,
	)


	logger.Info(
		"starting FaultPlane daemon",
	)


	// Root context with shutdown handling
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()



	// Telemetry initialization
	telemetryRegistry :=
		telemetry.NewRegistry()


	collector :=
		telemetry.NewCollector(
			telemetryRegistry,
		)



	// Storage / control layer
	controller :=
		control.New(
			logger,
			collector,
		)



	// Gateway layer
	gatewayManager :=
		gateway.New(
			logger,
		)



	// API layer
	router :=
		api.NewRouter(
			controller,
		)



	// Runtime daemon
	daemon :=
		runtime.New(
			logger,
			controller,
			gatewayManager,
			router,
		)



	// Start daemon
	go func(){

		if err := daemon.Start(ctx); err != nil {

			logger.Error(
				"daemon stopped with error",
				"error",
				err,
			)

			cancel()
		}

	}()



	// Wait shutdown signal
	<-ctx.Done()



	logger.Info(
		"shutdown signal received",
	)



	// Graceful shutdown timeout
	shutdownCtx, shutdownCancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

	defer shutdownCancel()



	if err :=
		daemon.Stop(shutdownCtx); err != nil {


		logger.Error(
			"graceful shutdown failed",
			"error",
			err,
		)

		os.Exit(1)
	}



	logger.Info(
		"FaultPlane daemon stopped cleanly",
	)
}