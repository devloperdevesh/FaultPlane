//go:build windows

package kernel

// Loader is a Windows-compatible no-op.
// FaultPlane's real eBPF loader is implemented in loader_linux.go.
type Loader struct {
	loaded bool
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(objectPath string) error {
	return nil
}

func (l *Loader) Loaded() bool {
	return l.loaded
}

func (l *Loader) Close() error {
	l.loaded = false
	return nil
}
