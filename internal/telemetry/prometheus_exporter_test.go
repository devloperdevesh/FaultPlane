package telemetry

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrometheusHandler(t *testing.T) {
	registry := NewRegistry()

	registry.RecordKernelEvent("kernel_warning", "link")
	registry.RecordKernelTransition("healthy", "degraded", "kernel_warning", "enforce", "kernel/rtnetlink")
	registry.RecordKernelError(errors.New("kernel failure"))
	registry.RecordRecovery()

	handler := PrometheusHandler(registry)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics/prometheus", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "text/plain; version=0.0.4" {
		t.Fatalf("unexpected content type: %q", got)
	}

	body := rec.Body.String()

	expectedMetrics := []string{
		"faultplane_kernel_events_total",
		"faultplane_kernel_transitions_total",
		"faultplane_kernel_errors_total",
		"faultplane_kernel_enforcement_total",
		"faultplane_kernel_recoveries_total",
		"faultplane_kernel_loader_loads_total",
		"faultplane_kernel_loader_closes_total",
		"faultplane_kernel_loader_failures_total",
		"faultplane_requests_total",
		"faultplane_recoveries_total",
		"faultplane_checkpoints_total",
		"faultplane_average_latency_seconds",
		"faultplane_cpu_usage",
		"faultplane_memory_usage",
		"faultplane_kernel_last_state",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("expected Prometheus metric %q in output", metric)
		}
	}

	if !strings.Contains(body, `faultplane_kernel_last_state{state="degraded"} 1`) {
		t.Errorf("expected last kernel state to be exported")
	}
}

func TestPrometheusHandlerMethodNotAllowed(t *testing.T) {
	registry := NewRegistry()
	handler := PrometheusHandler(registry)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/metrics/prometheus", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestPrometheusHandlerNilRegistry(t *testing.T) {
	handler := PrometheusHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics/prometheus", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}
}
