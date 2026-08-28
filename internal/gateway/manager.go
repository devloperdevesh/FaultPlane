package gateway

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/config"
	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/kernel"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Manager struct {
	logger         *slog.Logger
	registry       *telemetry.Registry
	collector      *telemetry.Collector
	workerStore    *api.WorkerStore
	telemetryStore *api.TelemetryStore
	topology       *control.TopologyController
	controller     *control.Controller
	config         config.Config
}

func New(
	logger *slog.Logger,
	registry *telemetry.Registry,
	collector *telemetry.Collector,
	controller *control.Controller,
) *Manager {
	return &Manager{
		logger:         logger,
		registry:       registry,
		collector:      collector,
		workerStore:    api.NewWorkerStore(),
		telemetryStore: api.NewTelemetryStore(),
		topology:       control.NewTopologyController(),
		controller:     controller,
		config:         config.Load(),
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

func (m *Manager) Topology() *control.TopologyController {
	return m.topology
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
	m.logger.Info(
		"gateway starting",
		"host", m.config.Host,
		"port", m.config.Port,
	)

	handler := NewRouter(
		m.registry,
		m.collector,
		m.workerStore,
		m.telemetryStore,
		monitor,
		m.topology,
		m.controller,
	)

	server := NewServer(handler, m.config)

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
