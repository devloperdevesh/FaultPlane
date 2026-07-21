package gateway

import (
	"context"
	"net/http"
	"time"
)

// Server represents AgentMesh HTTP gateway server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a production HTTP server.
func NewServer(handler http.Handler) *Server {

	return &Server{
		httpServer: &http.Server{

			Addr: ":8080",

			Handler: handler,

			ReadTimeout: 10 * time.Second,

			WriteTimeout: 10 * time.Second,

			IdleTimeout: 60 * time.Second,
		},
	}
}

// Start starts gateway server.
func (s *Server) Start() error {

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops server.
func (s *Server) Shutdown(
	ctx context.Context,
) error {

	return s.httpServer.Shutdown(ctx)
}
