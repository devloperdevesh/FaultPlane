package bpf

import (
	"testing"
)

func TestLoader_InvalidObject(t *testing.T) {

	loader := NewLoader()

	err := loader.Load(
		"invalid-path/sockmap.bpf.o",
	)

	if err == nil {
		t.Fatal(
			"expected error when loading invalid bpf object",
		)
	}
}

func TestLoader_CloseWithoutLoad(t *testing.T) {

	loader := NewLoader()

	err := loader.Close()

	if err != nil {
		t.Fatalf(
			"unexpected close error: %v",
			err,
		)
	}
}

func TestLoader_DoubleLoadProtection(t *testing.T) {

	loader := NewLoader()

	loader.loaded = true

	err := loader.Load(
		"sockmap.bpf.o",
	)

	if err == nil {
		t.Fatal(
			"expected error when loading already loaded bpf collection",
		)
	}
}
