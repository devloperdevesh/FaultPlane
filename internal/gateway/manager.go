package gateway

import (
	"context"
	"log/slog"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Manager struct {
	logger    *slog.Logger
	server    *Server
	registry  *telemetry.Registry
	collector *telemetry.Collector
}

func New(
	logger *slog.Logger,
	registry *telemetry.Registry,
	collector *telemetry.Collector,
) *Manager {
	return &Manager{
		logger:    logger,
		registry:  registry,
		collector: collector,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	m.logger.Info("gateway starting")

	handler := NewRouter(
		m.registry,
		m.collector,
	)

	m.server = NewServer(handler)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := m.server.Shutdown(shutdownCtx); err != nil {
			m.logger.Error(
				"gateway shutdown failed",
				"error", err,
			)
		}
	}()

	m.logger.Info(
		"gateway listening",
		"address", ":8080",
	)

	return m.server.Start()
}
