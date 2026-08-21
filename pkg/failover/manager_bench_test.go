// SPDX-License-Identifier: Apache-2.0

package failover

import (
	"testing"
)

func BenchmarkFailoverTransition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		store := NewStateStore(StateHealthy)
		manager := NewManager(store)

		if _, err := manager.HandleFailure("benchmark degradation"); err != nil {
			b.Fatal(err)
		}

		if _, err := manager.HandleFailure("benchmark failure"); err != nil {
			b.Fatal(err)
		}

		if _, err := manager.HandleFailure("benchmark fallback"); err != nil {
			b.Fatal(err)
		}
	}
}
