package bpf

import (
	"fmt"
	"sync"

	"github.com/cilium/ebpf"
)

// Loader manages lifecycle of FaultPlane eBPF programs.
type Loader struct {
	mu sync.Mutex

	collection *ebpf.Collection
	loaded     bool
}

// NewLoader creates a new BPF loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load loads and initializes an eBPF collection.
func (l *Loader) Load(objectPath string) error {

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.loaded {
		return fmt.Errorf(
			"bpf collection already loaded",
		)
	}

	spec, err := ebpf.LoadCollectionSpec(objectPath)
	if err != nil {
		return fmt.Errorf(
			"failed to load bpf object %q: %w",
			objectPath,
			err,
		)
	}

	collection, err := ebpf.NewCollection(spec)
	if err != nil {
		return fmt.Errorf(
			"failed to create bpf collection: %w",
			err,
		)
	}

	l.collection = collection
	l.loaded = true

	return nil
}

// Close releases all loaded eBPF resources.
func (l *Loader) Close() error {

	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.loaded {
		return nil
	}

	if l.collection != nil {
		l.collection.Close()
	}

	l.collection = nil
	l.loaded = false

	return nil
}
