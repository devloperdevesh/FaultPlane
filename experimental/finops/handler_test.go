package finops

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func TestHandlerReturnsLiveTelemetry(t *testing.T) {
	registry := telemetry.NewRegistry()

	registry.IncRequests()
	registry.IncWorkers()
	registry.RecordCheckpoint()
	registry.RecordRecovery()

	handler := NewHandler(registry)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/finops?cpu_hourly=2&memory_gb_hourly=1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()

	for _, field := range []string{
		`"requests":1`,
		`"workers":1`,
		`"recoveries":1`,
		`"checkpoints":1`,
		`"estimated_cost"`,
		`"observed_at"`,
	} {
		if !strings.Contains(body, field) {
			t.Fatalf("response missing %q: %s", field, body)
		}
	}
}

func TestHandlerRejectsInvalidRate(t *testing.T) {
	registry := telemetry.NewRegistry()
	handler := NewHandler(registry)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/finops?cpu_hourly=invalid",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandlerRejectsNegativeRate(t *testing.T) {
	registry := telemetry.NewRegistry()
	handler := NewHandler(registry)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/finops?cpu_hourly=-1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
