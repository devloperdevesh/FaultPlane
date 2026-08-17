package bpf

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

func (l *Loader) Load(
	objectPath string,
) error {

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.loaded {
		return fmt.Errorf(
			"bpf program already loaded",
		)
	}

	spec, err := ebpf.LoadCollectionSpec(objectPath)
	if err != nil {
		return fmt.Errorf(
			"load bpf spec: %w",
			err,
		)
	}

	collection, err := ebpf.NewCollection(spec)
	if err != nil {
		return fmt.Errorf(
			"create bpf collection: %w",
			err,
		)
	}

	l.collection = collection
	l.loaded = true

	return nil
}

func (l *Loader) Close() error {

	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.loaded {
		return nil
	}

	l.collection.Close()

	l.loaded = false

	return nil
}
