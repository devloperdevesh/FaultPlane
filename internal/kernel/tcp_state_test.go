package kernel

import "testing"

func TestTCPStateParser_DisconnectEvent(t *testing.T) {
	parser := NewTCPStateParser()

	event := TCPEvent{
		SourceIP:        "10.0.0.1",
		DestinationIP:   "10.0.0.2",
		SourcePort:      8080,
		DestinationPort: 9000,
		State:           TCPStateClosed,
	}

	result, err := parser.Parse(event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Type != "connection_failure" {
		t.Fatalf(
			"expected connection_failure, got %s",
			result.Type,
		)
	}
}

func TestTCPStateParser_InvalidState(t *testing.T) {
	parser := NewTCPStateParser()

	event := TCPEvent{}

	_, err := parser.Parse(event)

	if err == nil {
		t.Fatal("expected invalid tcp state error")
	}
}

func TestParseRawState(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  TCPState
	}{
		{
			name:  "closed state",
			input: "closed",
			want:  TCPStateClosed,
		},
	}

	parser := ParseRawState

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := parser(tt.input)

			if got != tt.want {
				t.Fatalf(
					"expected %s, got %s",
					tt.want,
					got,
				)
			}

		})
	}
}
