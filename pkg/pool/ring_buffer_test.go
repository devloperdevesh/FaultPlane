package pool

import "testing"

func TestRingBuffer(t *testing.T) {
	b := NewRingBuffer(2)

	for i := uint64(1); i <= 4; i++ {
		if err := b.EnqueueTokenEnvelope(i); err != nil {
			t.Fatalf("enqueue %d failed: %v", i, err)
		}
	}

	if err := b.EnqueueTokenEnvelope(5); err != ErrBufferSaturated {
		t.Fatalf("expected saturation, got %v", err)
	}

	for i := uint64(1); i <= 4; i++ {
		value, err := b.DequeueTokenEnvelope()

		if err != nil {
			t.Fatalf("dequeue failed: %v", err)
		}

		if value != i {
			t.Fatalf("expected %d, got %d", i, value)
		}
	}

	if _, err := b.DequeueTokenEnvelope(); err != ErrBufferExhausted {
		t.Fatalf("expected exhaustion, got %v", err)
	}
}
