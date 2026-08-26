package verifier

import "testing"

func TestSafetyVerifier(t *testing.T) {
	v := NewSafetyVerifier()

	if err := v.AssertStateInvariant(10, 100); err != nil {
		t.Fatalf("valid state rejected: %v", err)
	}

	if err := v.AssertStateInvariant(101, 100); err != ErrBoundaryViolation {
		t.Fatalf("expected boundary violation, got %v", err)
	}

	if err := v.AssertStateInvariant(0, 100); err != ErrBoundaryViolation {
		t.Fatalf("expected zero descriptor violation, got %v", err)
	}
}
