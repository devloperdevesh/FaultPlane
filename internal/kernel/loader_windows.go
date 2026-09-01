//go:build windows

package kernel

import "fmt"

// Loader is a Windows-compatible placeholder.
// FaultPlane's real eBPF loader is implemented in loader_linux.go.
type Loader struct {
	loaded bool
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(objectPath string) error {
	return fmt.Errorf("FaultPlane eBPF loader requires Linux; Windows is unsupported")
}

func (l *Loader) Loaded() bool {
	return l.loaded
}

func (l *Loader) Close() error {
	l.loaded = false
	return nil
}
