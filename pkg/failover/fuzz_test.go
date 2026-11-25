// SPDX-License-Identifier: Apache-2.0
package failover

import (
	"testing"
)

// FuzzStateStoreSwap verifies that arbitrary malformed byte inputs or transition states
// do not cause panics, memory corruption, or race conditions during atomic state swapping.
func FuzzStateStoreSwap(f *testing.F) {
	f.Add([]byte("healthy"), []byte("failed"))
	f.Add([]byte("active"), []byte("degraded"))
	f.Add([]byte(""), []byte("malformed_packet_data_stream_overflow_test"))
	f.Add([]byte{0x00, 0xff, 0xfe, 0x12}, []byte{0x7f, 0x00, 0x00, 0x01})

	store := NewStateStore(StateHealthy)

	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		stateA := ConnectionState(string(a))
		stateB := ConnectionState(string(b))

		store.Swap(stateA)
		if loaded := store.Load(); loaded != stateA {
			t.Fatalf("expected state %s, got %s", stateA, loaded)
		}

		store.Swap(stateB)
		if loaded := store.Load(); loaded != stateB {
			t.Fatalf("expected state %s, got %s", stateB, loaded)
		}
	})
}
