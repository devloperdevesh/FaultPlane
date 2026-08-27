package gateway

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/kernel"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Manager struct {
	logger         *slog.Logger
	registry       *telemetry.Registry
	collector      *telemetry.Collector
	workerStore    *api.WorkerStore
	telemetryStore *api.TelemetryStore
}

func New(
	logger *slog.Logger,
	registry *telemetry.Registry,
	collector *telemetry.Collector,
) *Manager {
	return &Manager{
		logger:         logger,
		registry:       registry,
		collector:      collector,
		workerStore:    api.NewWorkerStore(),
		telemetryStore: api.NewTelemetryStore(),
	}
}

func (m *Manager) Registry() *telemetry.Registry {
	return m.registry
}

func (m *Manager) WorkerStore() *api.WorkerStore {
	return m.workerStore
}

func (m *Manager) TelemetryStore() *api.TelemetryStore {
	return m.telemetryStore
}

func (m *Manager) Start(ctx context.Context) error {
	return m.start(ctx, nil)
}

func (m *Manager) StartWithKernel(
	ctx context.Context,
	monitor *kernel.Monitor,
) error {
	return m.start(ctx, monitor)
}

func (m *Manager) start(
	ctx context.Context,
	monitor *kernel.Monitor,
) error {
	m.logger.Info("gateway starting")

	handler := NewRouter(
		m.registry,
		m.collector,
		m.workerStore,
		m.telemetryStore,
		monitor,
	)

	server := NewServer(handler)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_ = server.Shutdown(shutdownCtx)
	}()

	err := server.Start()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
