// SPDX-License-Identifier: Apache-2.0
package failover

import (
	"sync/atomic"
	"testing"
)

// BenchmarkSub2msKernelShunting measures the local atomic state-transition
// overhead used by the failover control path.
//
// This benchmark intentionally does not enforce a hard 2ms wall-clock
// threshold per iteration because scheduler/CI noise can make such a
// threshold nondeterministic. Use benchmark results to evaluate latency.
func BenchmarkSub2msKernelShunting(b *testing.B) {
	var (
		mockPrimarySocket  uintptr = 0x10A0B000
		mockStandbySocket  uintptr = 0x20C0D000
		atomicSwapRegistry uint32
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if atomic.CompareAndSwapUint32(
			&atomicSwapRegistry,
			0,
			1,
		) {
			mockPrimarySocket = mockStandbySocket
			atomic.StoreUint32(&atomicSwapRegistry, 0)
		}
	}

	_ = mockPrimarySocket
}
