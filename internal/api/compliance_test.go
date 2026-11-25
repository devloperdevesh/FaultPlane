package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func TestComplianceLedgerHandler(t *testing.T) {
	registry := telemetry.NewRegistry()
	handler := ComplianceLedgerHandler(registry)

	req := httptest.NewRequest(http.MethodGet, "/api/compliance/ledger", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp ComplianceLedgerResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.LedgerID == "" {
		t.Errorf("expected non-empty ledger ID")
	}

	if resp.LedgerSignature == "" {
		t.Errorf("expected non-empty ledger signature")
	}

	if resp.Summary.UptimePercentage < 99.9 {
		t.Errorf("expected uptime >= 99.9, got %f", resp.Summary.UptimePercentage)
	}

	if resp.Summary.ContextLossEvents != 0 {
		t.Errorf("expected 0 context loss events, got %d", resp.Summary.ContextLossEvents)
	}

	if resp.Summary.CacheAlignmentBytes != 64 {
		t.Errorf("expected 64-byte cache alignment, got %d", resp.Summary.CacheAlignmentBytes)
	}

	if len(resp.Logs) == 0 {
		t.Errorf("expected historical audit logs, got 0")
	}
}

func TestSignatureVerification(t *testing.T) {
	now := time.Now().UTC()
	sig1 := GenerateLedgerSignature("TEST-01", now, 99.999)
	sig2 := GenerateLedgerSignature("TEST-01", now, 99.999)

	if sig1 != sig2 {
		t.Errorf("signatures should be deterministic for identical inputs")
	}

	sig3 := GenerateLedgerSignature("TEST-02", now, 99.999)
	if sig1 == sig3 {
		t.Errorf("signatures should differ for different ledger IDs")
	}
}
