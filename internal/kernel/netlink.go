package kernel

import (
	"context"
	"fmt"
	"sync"

	"log/slog"
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

	n.running = true

	n.mu.Unlock()


	n.logger.Info(
		"kernel netlink listener started",
	)


	go func() {

		defer func() {

			n.mu.Lock()
			n.running = false
			n.mu.Unlock()

			close(n.events)

			n.logger.Info(
				"kernel netlink listener stopped",
			)

		}()


		for {

			select {

			case <-ctx.Done():

				return


			default:

				// Placeholder:
				// Real Linux netlink socket read
				// will be implemented here.

			}

		}

	}()


	return nil
}


// Stop gracefully shuts down listener.
func (n *NetlinkListener) Stop() {

	n.mu.Lock()

	defer n.mu.Unlock()


	if !n.running {
		return
	}

	n.logger.Info(
		"netlink shutdown requested",
	)
}