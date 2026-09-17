//go:build linux

package kernel

import "testing"

func TestValidateEvent_KnownTopologyEvents(t *testing.T) {
	tests := []struct {
		name string
		typ  string
	}{
		{name: "link", typ: "link"},
		{name: "address", typ: "address"},
		{name: "route", typ: "route"},
		{name: "neighbor", typ: "neighbor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{
				Type: tt.typ,
				Data: []byte{1, 2, 3},
			}

			if err := ValidateEvent(event); err != nil {
				t.Fatalf("expected valid event, got error: %v", err)
			}

			if !IsTopologyEvent(event) {
				t.Fatalf("expected topology event for type %q", tt.typ)
			}
		})
	}
}

func TestValidateEvent_UnknownType(t *testing.T) {
	event := Event{
		Type: "unknown",
		Data: []byte{1, 2, 3},
	}

	if err := ValidateEvent(event); err == nil {
		t.Fatal("expected unknown event type to fail validation")
	}

	if IsTopologyEvent(event) {
		t.Fatal("unknown event must not be classified as topology event")
	}
}

func TestValidateEvent_EmptyType(t *testing.T) {
	event := Event{
		Data: []byte{1, 2, 3},
	}

	if err := ValidateEvent(event); err == nil {
		t.Fatal("expected empty event type to fail validation")
	}
}

func TestValidateEvent_EmptyPayloadAllowed(t *testing.T) {
	event := Event{
		Type: "link",
	}

	if err := ValidateEvent(event); err != nil {
		t.Fatalf("expected structurally valid event with empty payload: %v", err)
	}
}
