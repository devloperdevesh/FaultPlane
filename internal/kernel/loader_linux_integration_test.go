//go:build linux

package kernel

import (
	"os"
	"testing"
)

func TestRealLoaderLoadAttachClose(t *testing.T) {
	if os.Getenv("FAULTPLANE_EBPF_INTEGRATION") != "1" {
		t.Skip("real eBPF integration disabled")
	}

	object := os.Getenv("FAULTPLANE_BPF_OBJECT")
	if object == "" {
		t.Fatal("FAULTPLANE_BPF_OBJECT is required")
	}

	loader := NewLoader()

	if loader.Loaded() {
		t.Fatal("loader unexpectedly loaded before Load")
	}

	if err := loader.Load(object); err != nil {
		t.Fatalf("REAL LOADER LOAD FAILED: %v", err)
	}

	if !loader.Loaded() {
		t.Fatal("loader not marked loaded after Load")
	}

	mode, err := loader.EnforcementMode()
	if err != nil {
		_ = loader.Close()
		t.Fatalf("read enforcement mode failed: %v", err)
	}

	if mode != EnforcementAllow {
		_ = loader.Close()
		t.Fatalf("unexpected initial enforcement mode: got %d want %d", mode, EnforcementAllow)
	}

	if err := loader.SetEnforcementMode(EnforcementDrop); err != nil {
		_ = loader.Close()
		t.Fatalf("set drop mode failed: %v", err)
	}

	mode, err = loader.EnforcementMode()
	if err != nil {
		_ = loader.Close()
		t.Fatalf("verify drop mode failed: %v", err)
	}

	if mode != EnforcementDrop {
		_ = loader.Close()
		t.Fatalf("drop mode mismatch: got %d want %d", mode, EnforcementDrop)
	}

	if err := loader.SetEnforcementMode(EnforcementAllow); err != nil {
		_ = loader.Close()
		t.Fatalf("restore allow mode failed: %v", err)
	}

	if err := loader.Close(); err != nil {
		t.Fatalf("REAL LOADER CLOSE FAILED: %v", err)
	}

	if loader.Loaded() {
		t.Fatal("loader still marked loaded after Close")
	}

	if err := loader.Close(); err != nil {
		t.Fatalf("second Close failed: %v", err)
	}
}
