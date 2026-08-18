package failover

import "testing"

func TestFailoverManagerFailureFlow(
	t *testing.T,
) {

	store := NewStateStore(
		StateHealthy,
	)

	manager := NewManager(
		store,
	)

	event, err := manager.HandleFailure(
		"backend timeout",
	)

	if err != nil {
		t.Fatal(err)
	}

	if event.Current != StateDegraded {
		t.Fatalf(
			"expected degraded got %s",
			event.Current,
		)
	}

	event, err = manager.HandleFailure(
		"connection lost",
	)

	if err != nil {
		t.Fatal(err)
	}

	if event.Current != StateFailed {
		t.Fatalf(
			"expected failed got %s",
			event.Current,
		)
	}
}

func TestFailoverRecoveryFlow(
	t *testing.T,
) {

	store := NewStateStore(
		StateFallback,
	)

	manager := NewManager(
		store,
	)

	event, err := manager.Recover(
		"primary restored",
	)

	if err != nil {
		t.Fatal(err)
	}

	if event.Current != StateRecovering {
		t.Fatalf(
			"expected recovering got %s",
			event.Current,
		)
	}
}
