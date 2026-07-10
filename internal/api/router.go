package api

import (
	"net/http"

	"github.com/devloperdevesh/agentmesh/internal/control"
)

func NewRouter(
	controller *control.Controller,
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

	mux.HandleFunc(
		"/metrics",
		MetricsHandler,
	)

	return LoggingMiddleware(mux)
}