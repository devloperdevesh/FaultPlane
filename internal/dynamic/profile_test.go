package dynamic

import "testing"

func TestHotPatcher(t *testing.T) {
	p := NewHotPatcher()

	if !p.InjectLiveInstructionPatch(0x1000, 0x1234) {
		t.Fatal("valid patch rejected")
	}

	s := p.Snapshot()

	if s.OpcodeOffset != 0x1000 {
		t.Fatalf("unexpected address: %#x", s.OpcodeOffset)
	}

	if s.ExecutionSequence != 1 {
		t.Fatalf("expected sequence 1, got %d", s.ExecutionSequence)
	}
}
