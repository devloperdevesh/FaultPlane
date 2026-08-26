package hyperplane

import "testing"

func TestHyperPlaneControllerValidatesExecutionWindow(t *testing.T) {
	c, err := NewHyperPlaneController(0x1000, 0x2000)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.ExecuteAsynchronousShunt(0x1800, 0x2000); err != nil {
		t.Fatalf("expected valid buffer, got %v", err)
	}

	if err := c.ExecuteAsynchronousShunt(0x3000, 0x4000); err != ErrContextContention {
		t.Fatalf("expected execution-window rejection, got %v", err)
	}
}

func TestHyperPlaneControllerRejectsInvalidConfiguration(t *testing.T) {
	if _, err := NewHyperPlaneController(0, 0); err != ErrInvalidMemoryWindow {
		t.Fatalf("expected invalid memory window, got %v", err)
	}
}
