package gateway

import (
	"context"
	"log/slog"
	"time"
)

type Manager struct {
	logger *slog.Logger
	server *Server
}

func New(logger *slog.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	m.logger.Info("gateway starting")

	handler := NewRouter()

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
