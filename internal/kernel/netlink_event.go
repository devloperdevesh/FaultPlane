//go:build linux

package kernel

import "fmt"

// IsTopologyEvent reports whether an event belongs to a known
// rtnetlink topology/event family.
func IsTopologyEvent(event Event) bool {
	switch event.Type {
	case "link",
		"address",
		"route",
		"neighbor":
		return true
	default:
		return false
	}
}

// ValidateEvent verifies the minimum runtime contract for a
// decoded Linux rtnetlink event.
func ValidateEvent(event Event) error {
	if event.Type == "" {
		return fmt.Errorf("netlink event type is empty")
	}

	if !IsTopologyEvent(event) {
		return fmt.Errorf("unknown netlink event type: %q", event.Type)
	}

	return nil
}
