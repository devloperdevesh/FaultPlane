package gateway

import (
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func TelemetryMiddleware(
	collector *telemetry.Collector,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		collector.RecordRequest()
		next.ServeHTTP(w, r)
	})
}
