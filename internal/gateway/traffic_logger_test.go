package gateway

import (
	"os"
	"strings"
	"testing"
)

func TestProductionTrafficLoggerWritesAuditRecord(t *testing.T) {
	f, err := os.CreateTemp("", "faultplane-traffic-*.log")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	logger := NewTrafficLogger(path)
	logger.LogActiveEnterpriseTraffic(1024)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)

	if !strings.Contains(content, "packets=1") {
		t.Fatalf("expected packet counter in audit record: %s", content)
	}

	if !strings.Contains(content, "total_bytes=1024") {
		t.Fatalf("expected byte counter in audit record: %s", content)
	}
}
