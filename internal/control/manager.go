package control

import (
	"context"
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
) *Manager {
	return &Manager{
		logger:     logger,
		controller: NewController(store, collector),
	}
}

func (m *Manager) Controller() *Controller {
	return m.controller
}

func (m *Manager) Start(ctx context.Context) error {
	m.logger.Info(
		"control plane started",
	)

	<-ctx.Done()

	m.logger.Info(
		"control plane stopped",
	)

	return nil
}
