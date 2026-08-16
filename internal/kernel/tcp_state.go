package kernel

import (
	"fmt"
	"strings"
)

// TCPState represents Linux TCP connection states.
type TCPState string

const (
	TCPStateEstablished TCPState = "ESTABLISHED"
	TCPStateSynSent     TCPState = "SYN_SENT"
	TCPStateSynRecv     TCPState = "SYN_RECV"
	TCPStateFinWait1    TCPState = "FIN_WAIT1"
	TCPStateFinWait2    TCPState = "FIN_WAIT2"
	TCPStateTimeWait    TCPState = "TIME_WAIT"
	TCPStateClosed      TCPState = "CLOSED"
	TCPStateCloseWait   TCPState = "CLOSE_WAIT"
	TCPStateLastAck     TCPState = "LAST_ACK"
	TCPStateListen      TCPState = "LISTEN"
)

// TCPEvent represents raw TCP kernel event.
type TCPEvent struct {
	SourceIP string

	DestinationIP string

	SourcePort uint16

	DestinationPort uint16

	State TCPState
}

// FaultEvent is internal FaultPlane event.
// Gateway and recovery consume this.
type FaultEvent struct {
	Type string

	Source string

	Reason string

	Metadata map[string]string
}

// TCPStateParser converts kernel TCP events
// into FaultPlane events.
type TCPStateParser struct{}

// NewTCPStateParser creates parser.
func NewTCPStateParser() *TCPStateParser {

	return &TCPStateParser{}
}

// Parse converts TCP event into FaultPlane event.
func (p *TCPStateParser) Parse(
	event TCPEvent,
) (*FaultEvent, error) {

	if event.State == "" {

		return nil, fmt.Errorf(
			"invalid tcp state",
		)

	}

	switch event.State {

	case TCPStateClosed,
		TCPStateCloseWait,
		TCPStateLastAck,
		TCPStateFinWait1,
		TCPStateFinWait2:

		return &FaultEvent{

			Type: "connection_failure",

			Source: "kernel/tcp",

			Reason: string(event.State),

			Metadata: map[string]string{

				"source_ip": event.SourceIP,

				"destination_ip": event.DestinationIP,

				"source_port": fmt.Sprintf(
					"%d",
					event.SourcePort,
				),

				"destination_port": fmt.Sprintf(
					"%d",
					event.DestinationPort,
				),
			},
		}, nil

	default:

		return &FaultEvent{

			Type: "connection_state_change",

			Source: "kernel/tcp",

			Reason: string(event.State),

			Metadata: map[string]string{

				"source_ip": event.SourceIP,

				"destination_ip": event.DestinationIP,
			},
		}, nil

	}

}

// ParseRawState converts string state.
func ParseRawState(
	state string,
) TCPState {

	return TCPState(
		strings.ToUpper(
			state,
		),
	)

}
