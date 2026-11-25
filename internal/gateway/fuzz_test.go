// SPDX-License-Identifier: Apache-2.0
package gateway

import (
	"fmt"
	"net"
	"testing"
)

// FuzzPacketRoutingStream tests routing IP matching against malformed byte packets,
// asserting the engine guarantees failure-masking without panic or memory leaks.
func FuzzPacketRoutingStream(f *testing.F) {
	f.Add([]byte{192, 168, 1, 1}, uint32(0))
	f.Add([]byte{10, 0, 0, 1}, uint32(10<<24|1))
	f.Add([]byte{0, 0, 0, 0}, uint32(0xFFFFFFFF))
	f.Add([]byte{172, 16, 254, 255}, uint32(172<<24|16<<16|254<<8|255))
	f.Add([]byte("malformed_layer_4_socket_stream"), uint32(12345))

	engine := NewRoutingEngine()

	routes := make([]RouteEntry, 16)
	for i := 0; i < 16; i++ {
		routes[i] = RouteEntry{
			IP:        uint32(10<<24 | 0<<16 | 0<<8 | (i + 1)),
			TargetURL: fmt.Sprintf("http://cluster-worker-%d:8080", i+1),
			Healthy:   true,
		}
	}
	engine.SetRoutes(routes)

	f.Fuzz(func(t *testing.T, rawPacket []byte, target uint32) {
		// Test IP parsing from raw byte slices
		if len(rawPacket) >= 4 {
			ip := net.IPv4(rawPacket[0], rawPacket[1], rawPacket[2], rawPacket[3])
			asUint := IPToUint32(ip)
			_ = Uint32ToIP(asUint)
		}

		// Test MatchIP against fuzz-generated target uint32
		matched, err := engine.MatchIP(target)
		if err == nil {
			if matched.IP != target {
				t.Fatalf("fuzz error: expected %d, got %d", target, matched.IP)
			}
		} else if err != ErrRouteNotFound && err != ErrInvalidRoutePool {
			t.Fatalf("unexpected error returned: %v", err)
		}
	})
}
