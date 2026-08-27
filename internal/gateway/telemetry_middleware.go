package gateway

import (
	"net/http"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func TelemetryMiddleware(
	collector *telemetry.Collector,
	store *api.TelemetryStore,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		collector.RecordRequest()

		if store != nil {
			store.Add(api.TelemetryEvent{
				Type:      "request",
				Timestamp: time.Now().UTC(),
				Metadata: map[string]string{
					"method": r.Method,
					"path":   r.URL.Path,
				},
			})
		}

		next.ServeHTTP(w, r)

		if store != nil {
			store.Add(api.TelemetryEvent{
				Type:      "request_complete",
				Timestamp: time.Now().UTC(),
				Value:     float64(time.Since(start).Microseconds()) / 1000,
				Metadata: map[string]string{
					"method": r.Method,
					"path":   r.URL.Path,
				},
			})
		}
	})
}
