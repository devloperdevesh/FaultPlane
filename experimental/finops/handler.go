package finops

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

// Handler exposes live FinOps estimates derived from the runtime telemetry
// registry.
//
// Rates are supplied explicitly through query parameters:
//   - cpu_hourly
//   - memory_gb_hourly
//
// The response is an estimate based on observed process telemetry. It is not
// provider billing or an invoice.
type Handler struct {
	Registry *telemetry.Registry
}

// NewHandler creates a live FinOps HTTP handler.
func NewHandler(registry *telemetry.Registry) *Handler {
	return &Handler{
		Registry: registry,
	}
}

// ServeHTTP returns a live FinOps snapshot.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Registry == nil {
		http.Error(
			w,
			"finops telemetry registry unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	cpuRate, err := queryRate(r, "cpu_hourly")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memoryRate, err := queryRate(r, "memory_gb_hourly")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snapshot, err := FromRegistry(
		h.Registry,
		Rates{
			CPUHourly:      cpuRate,
			MemoryGBHourly: memoryRate,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(snapshot)
}

func queryRate(r *http.Request, name string) (float64, error) {
	value := r.URL.Query().Get(name)

	if value == "" {
		return 0, nil
	}

	rate, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, &rateError{
			name:  name,
			value: value,
		}
	}

	return rate, nil
}

type rateError struct {
	name  string
	value string
}

func (e *rateError) Error() string {
	return "invalid FinOps rate " + e.name + "=" + e.value
}
