package gateway

import (
	"net/http"
)

// Router owns HTTP route registration.
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates gateway router.
func NewRouter() http.Handler {

	r := &Router{
		mux: http.NewServeMux(),
	}

	r.registerRoutes()

	return r
}

// registerRoutes defines public API surface.
func (r *Router) registerRoutes() {

	r.mux.HandleFunc(
		"/health",
		healthHandler,
	)

	r.mux.HandleFunc(
		"/metrics",
		metricsHandler,
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
