package runtime

import (
	"runtime"
	"testing"
	"time"
)

func TestProcessCPUSamplerDetectsCPUActivity(t *testing.T) {
	sampler := NewProcessCPUSampler()

	_ = sampler.Sample()

	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer close(done)

		var x uint64
		for {
			select {
			case <-stop:
				return
			default:
				x++
				x *= 3
				if x == 0 {
					runtime.Gosched()
				}
			}
		}
	}()

	time.Sleep(3 * time.Second)

	sample := sampler.Sample()

	close(stop)
	<-done

	if sample <= 0 {
		t.Fatalf("expected CPU activity > 0, got %.6f%%", sample)
	}

	if sample > 100 {
		t.Fatalf("CPU sample out of range: %.6f%%", sample)
	}

	t.Logf("CPU activity detected: %.6f%%", sample)
}
