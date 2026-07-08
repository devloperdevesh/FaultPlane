package api

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", HealthHandler)
	mux.HandleFunc("/api/v1/checkpoint", CheckpointHandler)
	mux.HandleFunc("/api/v1/recover", RecoverHandler)

	return LoggingMiddleware(mux)
}
