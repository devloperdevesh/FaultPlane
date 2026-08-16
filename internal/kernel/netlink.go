package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Event represents a kernel network event.
// Future versions can map this to tcp_set_state/eBPF events.
type Event struct {
	Type string
	Data []byte
}

// NetlinkListener provides an abstraction over kernel event streaming.
type NetlinkListener struct {
	logger *slog.Logger

	events chan Event

	mu      sync.Mutex
	running bool

	cancel context.CancelFunc
}

// NewNetlinkListener creates a new kernel event listener.
func NewNetlinkListener(
	logger *slog.Logger,
) *NetlinkListener {

	return &NetlinkListener{
		logger: logger,
		events: make(chan Event, 256),
	}
}

// Events returns read-only kernel events stream.
func (n *NetlinkListener) Events() <-chan Event {
	return n.events
}

// Start initializes the kernel listener lifecycle.
func (n *NetlinkListener) Start(
	ctx context.Context,
) error {

	n.mu.Lock()

	if n.running {
		n.mu.Unlock()

		return fmt.Errorf(
			"netlink listener already running",
		)
	}

	listenerCtx, cancel := context.WithCancel(ctx)

	n.cancel = cancel
	n.running = true

	n.mu.Unlock()

	n.logger.Info(
		"kernel netlink listener started",
	)

	go n.loop(listenerCtx)

	return nil
}

// loop handles kernel events.
func (n *NetlinkListener) loop(
	ctx context.Context,
) {

	defer func() {

		n.mu.Lock()
		n.running = false
		n.mu.Unlock()

		n.logger.Info(
			"kernel netlink listener stopped",
		)

	}()

	ticker := time.NewTicker(
		100 * time.Millisecond,
	)

	defer ticker.Stop()

	for {

		select {

		case <-ctx.Done():

			return

		case <-ticker.C:

			// TODO:
			// Replace this ticker with real Linux
			// netlink socket receive.
			//
			// Future:
			// tcp_set_state
			// sock_diag
			// eBPF events

		}

	}

}

// Stop gracefully shuts down listener.
func (n *NetlinkListener) Stop() {

	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.running {
		return
	}

	if n.cancel != nil {

		n.cancel()
	}

	n.logger.Info(
		"netlink shutdown requested",
	)

}
