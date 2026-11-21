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
	"github.com/devloperdevesh/FaultPlane/internal/smart"
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
	smartRuntime   *smart.Runtime
	routingEngine  *RoutingEngine
	config         config.Config
}

func New(
	logger *slog.Logger,
	registry *telemetry.Registry,
	collector *telemetry.Collector,
	controller *control.Controller,
	smartRuntime *smart.Runtime,
) *Manager {
	events := api.NewTelemetryStore()
	engine := NewRoutingEngine()
	return &Manager{
		logger:         logger,
		registry:       registry,
		collector:      collector,
		workerStore:    api.NewWorkerStoreWithTelemetry(events),
		telemetryStore: events,
		topology:       control.NewTopologyController(),
		controller:     controller,
		smartRuntime:   smartRuntime,
		routingEngine:  engine,
		config:         config.Load(),
	}
}

func (m *Manager) RoutingEngine() *RoutingEngine {
	return m.routingEngine
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

	m.logger.Info(
		"hardware vector registry initialized",
		"track", m.routingEngine.ActiveTrack(),
		"has_avx512", m.routingEngine.Profile().HasAVX512F,
		"vector_bits", m.routingEngine.Profile().VectorBits,
		"parallel_lanes", m.routingEngine.Profile().ParallelLanes,
	)

	handler := NewRouterWithEngine(
		m.registry,
		m.collector,
		m.workerStore,
		m.telemetryStore,
		monitor,
		m.topology,
		m.controller,
		m.smartRuntime,
		m.routingEngine,
	)

	server := NewServer(handler, m.config)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			m.logger.Error(
				"gateway shutdown failed",
				"error", err,
			)
		}
	}()

	err := server.Start()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
