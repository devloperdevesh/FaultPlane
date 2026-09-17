//go:build linux

package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	rtnetlink "github.com/jsimonetti/rtnetlink/v2"

	"github.com/mdlayher/netlink"
)

const netlinkEventGroups uint32 = 0x00000055

const (
	netlinkEventBuffer  = 256
	netlinkReceiveRetry = 100 * time.Millisecond
)

// Event represents a kernel network event.
//
// Type identifies the rtnetlink message family/event.
// Data contains the raw serialized rtnetlink payload when available.
type Event struct {
	Type string
	Data []byte
}

// NetlinkListener provides a real Linux rtnetlink event stream.
type NetlinkListener struct {
	logger *slog.Logger

	mu      sync.Mutex
	running bool
	cancel  context.CancelFunc
	conn    *rtnetlink.Conn

	events chan Event
	done   chan struct{}
}

// NewNetlinkListener creates a Linux rtnetlink listener.
func NewNetlinkListener(logger *slog.Logger) *NetlinkListener {
	if logger == nil {
		logger = slog.Default()
	}

	return &NetlinkListener{
		logger: logger,
		events: make(chan Event, netlinkEventBuffer),
	}
}

// Events returns the read-only kernel event stream.
func (n *NetlinkListener) Events() <-chan Event {
	return n.events
}

// Start opens a real AF_NETLINK/NETLINK_ROUTE socket and begins
// receiving kernel rtnetlink messages.
func (n *NetlinkListener) Start(ctx context.Context) error {
	n.mu.Lock()

	if n.running {
		n.mu.Unlock()
		return fmt.Errorf("netlink listener already running")
	}

	conn, err := rtnetlink.Dial(&netlink.Config{
		Groups: netlinkEventGroups,
	})
	if err != nil {
		n.mu.Unlock()
		return fmt.Errorf("dial rtnetlink: %w", err)
	}

	listenerCtx, cancel := context.WithCancel(ctx)

	n.conn = conn
	n.cancel = cancel
	n.running = true
	n.done = make(chan struct{})

	done := n.done

	n.mu.Unlock()

	n.logger.Info("kernel rtnetlink listener started")

	go n.loop(listenerCtx, conn, done)

	return nil
}

// loop receives actual rtnetlink messages from the Linux kernel.
func (n *NetlinkListener) loop(
	ctx context.Context,
	conn *rtnetlink.Conn,
	done chan struct{},
) {
	defer func() {
		_ = conn.Close()

		n.mu.Lock()
		if n.conn == conn {
			n.conn = nil
			n.cancel = nil
			n.running = false
		}
		n.mu.Unlock()

		close(done)

		n.logger.Info("kernel rtnetlink listener stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msgs, rawMsgs, err := conn.Receive()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n.logger.Error(
				"rtnetlink receive failed",
				"error", err,
			)

			select {
			case <-ctx.Done():
				return
			case <-time.After(netlinkReceiveRetry):
			}

			continue
		}

		for i, msg := range msgs {
			event := Event{
				Type: rtnetlinkMessageType(msg),
			}

			if i < len(rawMsgs) {
				event.Data = append([]byte(nil), rawMsgs[i].Data...)
			}

			select {
			case n.events <- event:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Stop gracefully shuts down the rtnetlink listener.
func (n *NetlinkListener) Stop() {
	n.mu.Lock()

	if !n.running {
		n.mu.Unlock()
		return
	}

	cancel := n.cancel
	done := n.done

	n.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	// Closing the socket unblocks a Receive call immediately.
	n.mu.Lock()
	conn := n.conn
	n.mu.Unlock()

	if conn != nil {
		_ = conn.Close()
	}

	if done != nil {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			n.logger.Warn("timed out waiting for rtnetlink listener shutdown")
		}
	}
}

func rtnetlinkMessageType(msg rtnetlink.Message) string {
	switch msg.(type) {
	case *rtnetlink.LinkMessage:
		return "link"
	case *rtnetlink.AddressMessage:
		return "address"
	case *rtnetlink.RouteMessage:
		return "route"
	case *rtnetlink.NeighMessage:
		return "neighbor"
	default:
		return "unknown"
	}
}
