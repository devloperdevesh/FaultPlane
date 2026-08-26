package controlplane

import "testing"

func TestPolicyEngineEnforceSovereignTopology(t *testing.T) {
	e, err := NewPolicyEngine(4, 1000)
	if err != nil {
		t.Fatal(err)
	}

	if err := e.EnforceSovereignTopology(4, 500); err != nil {
		t.Fatalf("expected valid policy, got %v", err)
	}

	if err := e.EnforceSovereignTopology(5, 500); err != ErrTopologyMismatched {
		t.Fatalf("expected topology mismatch, got %v", err)
	}

	if err := e.EnforceSovereignTopology(4, 1001); err != ErrBandwidthSaturated {
		t.Fatalf("expected bandwidth saturation, got %v", err)
	}
}

func TestPolicyEngineRejectsInvalidConfiguration(t *testing.T) {
	if _, err := NewPolicyEngine(4, 0); err != ErrInvalidPolicy {
		t.Fatalf("expected invalid policy, got %v", err)
	}
}
