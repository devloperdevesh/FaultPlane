package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/config"
)

// Server represents the FaultPlane HTTP gateway server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a production HTTP server using the runtime configuration.
func NewServer(
	handler http.Handler,
	cfg config.Config,
) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: cfg.Host + ":" + cfg.Port,

			Handler: handler,

			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start starts the gateway HTTP server.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
