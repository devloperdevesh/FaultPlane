package gateway

import (
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/kernel"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(
	registry *telemetry.Registry,
	collector *telemetry.Collector,
	workerStore *api.WorkerStore,
	telemetryStore *api.TelemetryStore,
	kernelMonitor *kernel.Monitor,
) http.Handler {
	r := &Router{
		mux: http.NewServeMux(),
	}

	r.registerRoutes(
		registry,
		workerStore,
		telemetryStore,
		kernelMonitor,
	)

	return TelemetryMiddleware(
		collector,
		telemetryStore,
		r,
	)
}

func (r *Router) registerRoutes(
	registry *telemetry.Registry,
	workerStore *api.WorkerStore,
	telemetryStore *api.TelemetryStore,
	kernelMonitor *kernel.Monitor,
) {
	r.mux.HandleFunc(
		"/health",
		healthHandler,
	)

	r.mux.Handle(
		"/api/metrics",
		metricsHandler(registry),
	)

	r.mux.Handle(
		"/api/workers",
		api.WorkersHandler(workerStore),
	)

	r.mux.Handle(
		"/api/workers/",
		api.WorkersHandler(workerStore),
	)

	r.mux.Handle(
		"/api/telemetry",
		api.TelemetryHandler(telemetryStore),
	)

	r.mux.Handle(
		"/api/logs",
		api.LogsHandler(telemetryStore),
	)

	r.mux.Handle(
		"/api/network/events",
		api.NetworkEventsHandler(telemetryStore),
	)

	if kernelMonitor != nil {
		ebpfHandler := api.NewEBPFHandler(kernelMonitor)

		r.mux.Handle(
			"/api/ebpf/status",
			http.HandlerFunc(ebpfHandler.Status),
		)

		r.mux.Handle(
			"/api/ebpf/events",
			http.HandlerFunc(ebpfHandler.Events),
		)

		r.mux.Handle(
			"/api/ebpf/hooks",
			http.HandlerFunc(ebpfHandler.Hooks),
		)
	}
}

func (r *Router) ServeHTTP(
	w http.ResponseWriter,
	req *http.Request,
) {
	r.mux.ServeHTTP(w, req)
}
