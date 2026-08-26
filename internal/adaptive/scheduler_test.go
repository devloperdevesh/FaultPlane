package adaptive

import (
	"testing"
	"time"
)

func TestAdaptiveScheduler(t *testing.T) {
	s := NewAdaptiveScheduler(2 * time.Millisecond)

	if !s.EvaluateProviderSoundness(1 * time.Millisecond) {
		t.Fatal("healthy provider was rejected")
	}

	if s.EvaluateProviderSoundness(3 * time.Millisecond) {
		t.Fatal("slow provider was accepted")
	}

	snapshot := s.Snapshot()

	if snapshot.ProviderDropFlag != 1 {
		t.Fatalf("expected provider drop flag, got %d", snapshot.ProviderDropFlag)
	}
}
