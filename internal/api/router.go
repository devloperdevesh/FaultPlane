package api

import (
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func NewRouter(
	controller *control.Controller,
	registry *telemetry.Registry,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/health",
		HealthHandler,
	)

	mux.HandleFunc(
		"/checkpoint",
		CheckpointHandler(controller),
	)

	mux.HandleFunc(
		"/recover",
		RecoverHandler(controller),
	)

	mux.Handle(
		"/metrics",
		MetricsHandler(registry),
	)

	return LoggingMiddleware(mux)
}
