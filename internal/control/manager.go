package control

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/storage"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Manager struct {
	logger     *slog.Logger
	controller *Controller
}

func New(
	logger *slog.Logger,
	store storage.Store,
	collector *telemetry.Collector,
	runtime RuntimeExecutor,
) *Manager {
	if logger == nil {
		logger = slog.Default()
	}

	return &Manager{
		logger:     logger,
		controller: NewController(store, collector, runtime),
	}
}

func (m *Manager) Controller() *Controller {
	return m.controller
}

// SetRuntime wires the live runtime implementation into the control plane.
// The daemon calls this after constructing the platform-specific runtime
// and before starting the control manager.
func (m *Manager) SetRuntime(runtime RuntimeExecutor) error {
	if m == nil || m.controller == nil {
		return fmt.Errorf("set runtime: control manager is not initialized")
	}

	return m.controller.SetRuntime(runtime)
}

func (m *Manager) Start(ctx context.Context) error {
	if ctx == nil {
		return context.Canceled
	}

	m.logger.Info("control plane started")

	<-ctx.Done()

	m.logger.Info("control plane stopped")
	return nil
}
