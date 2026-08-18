package failover

import "testing"

func TestAtomicStateSwap(t *testing.T) {

	store := NewStateStore(
		StateHealthy,
	)

	if got := store.Load(); got != StateHealthy {

		t.Fatalf(
			"expected healthy state, got %s",
			got,
		)
	}

	store.Swap(
		StateFailed,
	)

	if got := store.Load(); got != StateFailed {

		t.Fatalf(
			"expected failed state, got %s",
			got,
		)
	}
}

func BenchmarkAtomicSwap(
	b *testing.B,
) {

	store := NewStateStore(
		StateHealthy,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		store.Swap(
			StateFailed,
		)

		store.Swap(
			StateHealthy,
		)
	}
}
