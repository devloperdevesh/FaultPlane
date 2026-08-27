package kernel

import (
	"fmt"
	"sync"

	"github.com/cilium/ebpf"
)

type Loader struct {
	mu sync.Mutex

	collection *ebpf.Collection
	loaded     bool
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(objectPath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.loaded {
		return fmt.Errorf("eBPF collection already loaded")
	}

	spec, err := ebpf.LoadCollectionSpec(objectPath)
	if err != nil {
		return fmt.Errorf("load eBPF spec: %w", err)
	}

	collection, err := ebpf.NewCollection(spec)
	if err != nil {
		return fmt.Errorf("create eBPF collection: %w", err)
	}

	l.collection = collection
	l.loaded = true

	return nil
}

func (l *Loader) Loaded() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.loaded
}

func (l *Loader) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.loaded {
		return nil
	}

	if l.collection != nil {
		l.collection.Close()
		l.collection = nil
	}

	l.loaded = false

	return nil
}
