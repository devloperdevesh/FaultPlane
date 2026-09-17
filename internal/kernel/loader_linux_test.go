package kernel

import "testing"

func TestLoaderInitialState(t *testing.T) {
	loader := NewLoader()

	if loader == nil {
		t.Fatal("NewLoader() returned nil")
	}
	if loader.Loaded() {
		t.Fatal("new loader must not be loaded")
	}
	if err := loader.Close(); err != nil {
		t.Fatalf("Close() on unloaded loader failed: %v", err)
	}
	if loader.Loaded() {
		t.Fatal("loader became loaded after Close()")
	}
}

func TestLoaderCloseIsIdempotent(t *testing.T) {
	loader := NewLoader()

	for i := 0; i < 3; i++ {
		if err := loader.Close(); err != nil {
			t.Fatalf("Close() call %d failed: %v", i+1, err)
		}
		if loader.Loaded() {
			t.Fatalf("loader marked loaded after Close() call %d", i+1)
		}
	}
}

func TestLoaderLoadInvalidObject(t *testing.T) {
	loader := NewLoader()

	err := loader.Load("/definitely/not/a/faultplane-ebpf-object.o")
	if err == nil {
		t.Fatal("Load() with invalid object path unexpectedly succeeded")
	}
	if loader.Loaded() {
		t.Fatal("loader marked loaded after failed Load()")
	}
	if closeErr := loader.Close(); closeErr != nil {
		t.Fatalf("Close() after failed Load() failed: %v", closeErr)
	}
}

func TestLoaderRejectsRepeatedLoadState(t *testing.T) {
	loader := NewLoader()

	loader.mu.Lock()
	loader.loaded = true
	loader.mu.Unlock()

	err := loader.Load("/definitely/not/a/faultplane-ebpf-object.o")
	if err == nil {
		t.Fatal("Load() unexpectedly succeeded while loader was already loaded")
	}
	if err.Error() != "eBPF collection already loaded" {
		t.Fatalf("unexpected repeated Load() error: %v", err)
	}

	loader.mu.Lock()
	loader.loaded = false
	loader.mu.Unlock()
}

func TestLoaderEnforcementRequiresLoaded(t *testing.T) {
	loader := NewLoader()

	if err := loader.SetEnforcementMode(EnforcementAllow); err == nil {
		t.Fatal("SetEnforcementMode() succeeded while loader was unloaded")
	}
	if _, err := loader.EnforcementMode(); err == nil {
		t.Fatal("EnforcementMode() succeeded while loader was unloaded")
	}
}

func TestLoaderRejectsInvalidEnforcementMode(t *testing.T) {
	loader := NewLoader()

	err := loader.SetEnforcementMode(EnforcementDrop + 1)
	if err == nil {
		t.Fatal("invalid enforcement mode was accepted")
	}
}
