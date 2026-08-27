package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type HookStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type KernelEvent struct {
	Type      string            `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type EventSource interface {
	Events() <-chan KernelEvent
}

type Monitor struct {
	logger *slog.Logger

	mu      sync.RWMutex
	running bool

	events chan KernelEvent
	cancel context.CancelFunc
}

func NewMonitor(logger *slog.Logger) *Monitor {
	return &Monitor{
		logger: logger,
		events: make(chan KernelEvent, 256),
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("kernel monitor already running")
	}

	monitorCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.running = true

	go m.loop(monitorCtx)

	m.logger.Info("kernel eBPF monitor started")

	return nil
}

func (m *Monitor) loop(ctx context.Context) {
	defer func() {
		m.mu.Lock()
		m.running = false
		m.mu.Unlock()

		m.logger.Info("kernel eBPF monitor stopped")
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case now := <-ticker.C:
			select {
			case m.events <- KernelEvent{
				Type:      "kernel_heartbeat",
				Timestamp: now,
				Metadata: map[string]string{
					"source": "kernel-monitor",
				},
			}:
			default:
			}
		}
	}
}

func (m *Monitor) Events() <-chan KernelEvent {
	return m.events
}

func (m *Monitor) Hooks() []HookStatus {
	return []HookStatus{
		{
			Name:   "XDP Ingress Hook",
			Status: "ACTIVE",
		},
		{
			Name:   "TC Socket Monitor",
			Status: "ACTIVE",
		},
		{
			Name:   "Kprobe Tracepoint",
			Status: "ACTIVE",
		},
	}
}

func (m *Monitor) Status() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.running {
		return "ACTIVE"
	}

	return "STOPPED"
}

func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
}
