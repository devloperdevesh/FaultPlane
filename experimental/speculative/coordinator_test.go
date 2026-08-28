package speculative

import (
	"errors"
	"sync"
	"testing"
)

func TestCoordinatorCommit(t *testing.T) {
	c := NewCoordinator()

	err := c.Register(Candidate{
		ID:      "candidate-a",
		Step:    42,
		Payload: []byte("real-state"),
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	got, err := c.Commit("candidate-a")
	if err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	if got.Status != StatusCommitted {
		t.Fatalf("expected committed status, got %q", got.Status)
	}

	if string(got.Payload) != "real-state" {
		t.Fatalf("payload mismatch: %q", got.Payload)
	}

	committed, ok := c.Committed()
	if !ok {
		t.Fatal("expected committed candidate")
	}

	if committed.ID != "candidate-a" {
		t.Fatalf("unexpected committed candidate: %q", committed.ID)
	}
}

func TestCoordinatorAllowsOnlyOneCommit(t *testing.T) {
	c := NewCoordinator()

	for _, id := range []string{"candidate-a", "candidate-b"} {
		if err := c.Register(Candidate{ID: id}); err != nil {
			t.Fatalf("register %s failed: %v", id, err)
		}
	}

	if _, err := c.Commit("candidate-a"); err != nil {
		t.Fatalf("first commit failed: %v", err)
	}

	if _, err := c.Commit("candidate-b"); err == nil {
		t.Fatal("expected second commit to fail")
	}
}

func TestCoordinatorCancel(t *testing.T) {
	c := NewCoordinator()

	if err := c.Register(Candidate{
		ID:      "candidate-a",
		Payload: []byte("discard-me"),
	}); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := c.Cancel("candidate-a"); err != nil {
		t.Fatalf("cancel failed: %v", err)
	}

	got, ok := c.Get("candidate-a")
	if !ok {
		t.Fatal("candidate disappeared after cancellation")
	}

	if got.Status != StatusCancelled {
		t.Fatalf("expected cancelled status, got %q", got.Status)
	}

	if len(got.Payload) != 0 {
		t.Fatal("cancelled candidate retained payload")
	}

	if _, err := c.Commit("candidate-a"); !errors.Is(err, ErrAlreadyCancelled) {
		t.Fatalf("expected ErrAlreadyCancelled, got %v", err)
	}
}

func TestCoordinatorConcurrentRegistrationAndCommit(t *testing.T) {
	c := NewCoordinator()

	const candidates = 32

	var wg sync.WaitGroup
	wg.Add(candidates)

	for i := 0; i < candidates; i++ {
		go func(i int) {
			defer wg.Done()

			id := "candidate-" + string(rune('a'+i))

			if err := c.Register(Candidate{
				ID:      id,
				Step:    uint64(i),
				Payload: []byte(id),
			}); err != nil {
				t.Errorf("register %s failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	if _, err := c.Commit("candidate-a"); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	committed, ok := c.Committed()
	if !ok {
		t.Fatal("expected committed candidate")
	}

	if committed.ID != "candidate-a" {
		t.Fatalf("unexpected committed candidate: %q", committed.ID)
	}
}
