package gateway

import (
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

// Router owns HTTP route registration.
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates gateway router.
func NewRouter(
	registry *telemetry.Registry,
	collector *telemetry.Collector,
) http.Handler {
	r := &Router{
		mux: http.NewServeMux(),
	}

	r.registerRoutes(registry)

	return TelemetryMiddleware(
		collector,
		r,
	)
}

// registerRoutes defines gateway public API surface.
func (r *Router) registerRoutes(
	registry *telemetry.Registry,
) {
	r.mux.HandleFunc(
		"/health",
		healthHandler,
	)

	r.mux.Handle(
		"/api/metrics",
		metricsHandler(registry),
	)

	r.mux.HandleFunc(
		"/workers",
		workersHandler,
	)
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(
	w http.ResponseWriter,
	req *http.Request,
) {
	r.mux.ServeHTTP(w, req)
}
