// SPDX-License-Identifier: Apache-2.0
package bpf

import (
	"testing"
)

func TestVectorBridgeStubs(t *testing.T) {
	ips := make([]uint32, 16)
	for i := 0; i < 16; i++ {
		ips[i] = uint32(192<<24 | 168<<16 | 1<<8 | (i + 1))
	}

	target := uint32(192<<24 | 168<<16 | 1<<8 | 5)

	if CgoEnabled && IsAVX512Compiled() {
		mask := NativeMatch16IPs(ips, target)
		expectedBit := uint16(1 << 4) // index 4 has target
		if mask&expectedBit == 0 {
			t.Fatalf("expected mask bit %x set, got %x", expectedBit, mask)
		}

		idx := NativeEvaluateIPs(ips, target)
		if idx != 4 {
			t.Fatalf("expected matched index 4, got %d", idx)
		}
	} else if !CgoEnabled {
		// In pure Go / non-Cgo fallback mode, ensure stubs are non-blocking and safe
		mask := NativeMatch16IPs(ips, target)
		if mask != 0 {
			t.Fatalf("expected 0 from non-cgo stub, got %d", mask)
		}

		idx := NativeEvaluateIPs(ips, target)
		if idx != -1 {
			t.Fatalf("expected -1 from non-cgo stub, got %d", idx)
		}
	} else {
		// Cgo is available but the runner CPU lacks AVX-512: native functions
		// must use a scalar fallback instead of executing AVX-512 instructions.
		mask := NativeMatch16IPs(ips, target)
		expectedBit := uint16(1 << 4)
		if mask&expectedBit == 0 {
			t.Fatalf("expected scalar mask bit %x set, got %x", expectedBit, mask)
		}

		idx := NativeEvaluateIPs(ips, target)
		if idx != 4 {
			t.Fatalf("expected scalar matched index 4, got %d", idx)
		}
	}
}
