package control

import (
	"context"
	"log/slog"
)

type Manager struct {
	logger *slog.Logger
}


func New(logger *slog.Logger) *Manager {

	return &Manager{
		logger: logger,
	}
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