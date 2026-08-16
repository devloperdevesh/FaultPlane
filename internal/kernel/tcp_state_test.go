package kernel

import "testing"

func TestTCPStateParserFailure(t *testing.T) {

	parser := NewTCPStateParser()

	event := TCPEvent{

		SourceIP: "10.0.0.1",

		DestinationIP: "10.0.0.2",

		SourcePort: 8080,

		DestinationPort: 9000,

		State: TCPStateClosed,
	}

	result, err := parser.Parse(event)

	if err != nil {
		t.Fatal(err)
	}

	if result.Type != "connection_failure" {

		t.Fatalf(
			"expected failure event got %s",
			result.Type,
		)

	}

}
