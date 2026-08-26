package tensorplane

import "testing"

func TestTensorIngressValidatesBuffer(t *testing.T) {
	c, err := NewTensorController(0x1000, 8)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.SplicedLineRateIngress(0x1100, 0x2000); err != nil {
		t.Fatalf("expected valid buffer, got %v", err)
	}

	if err := c.SplicedLineRateIngress(0, 0x2000); err != ErrBarAddressExceeded {
		t.Fatalf("expected address rejection, got %v", err)
	}

	if err := c.SplicedLineRateIngress(0x3000, 0x2000); err != ErrBarAddressExceeded {
		t.Fatalf("expected upper-bound rejection, got %v", err)
	}
}

func TestTensorControllerRejectsInvalidRingSize(t *testing.T) {
	if _, err := NewTensorController(0x1000, 32); err != ErrInvalidRingSize {
		t.Fatalf("expected invalid ring size, got %v", err)
	}
}
