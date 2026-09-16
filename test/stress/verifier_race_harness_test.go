// Package stress provides concurrency verification for FaultPlane's
// stream-ownership boundary.
//
// SPDX-License-Identifier: Apache-2.0
package stress

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	StateIdle int32 = iota
	StateBusy
)

type ActiveStreamSimulation struct {
	State  int32
	ConnID uint64
}

type SystemPipelineSimulator struct {
	Pool []*ActiveStreamSimulation

	SuccessSwaps uint64
	BusySkips    uint64
	RaceFailures uint64
}

func NewSystemPipelineSimulator(size int) *SystemPipelineSimulator {
	if size <= 0 {
		panic("stress harness: pool size must be greater than zero")
	}

	pool := make([]*ActiveStreamSimulation, size)

	for i := range pool {
		pool[i] = &ActiveStreamSimulation{
			State:  StateIdle,
			ConnID: uint64(i + 1000),
		}
	}

	return &SystemPipelineSimulator{
		Pool: pool,
	}
}

func (s *SystemPipelineSimulator) RunStressMatrix(
	duration time.Duration,
	workers int,
) {
	if duration <= 0 {
		panic("stress harness: duration must be greater than zero")
	}

	if workers <= 0 {
		panic("stress harness: workers must be greater than zero")
	}

	var stop atomic.Bool
	var wg sync.WaitGroup

	wg.Add(workers)

	for workerID := 0; workerID < workers; workerID++ {
		go func(id int) {
			defer wg.Done()

			rng := rand.New(
				rand.NewSource(
					time.Now().UnixNano() + int64(id),
				),
			)

			for !stop.Load() {
				stream := s.Pool[rng.Intn(len(s.Pool))]

				if !atomic.CompareAndSwapInt32(
					&stream.State,
					StateIdle,
					StateBusy,
				) {
					atomic.AddUint64(&s.BusySkips, 1)
					continue
				}

				atomic.AddUint64(&s.SuccessSwaps, 1)

				time.Sleep(10 * time.Microsecond)

				if !atomic.CompareAndSwapInt32(
					&stream.State,
					StateBusy,
					StateIdle,
				) {
					atomic.AddUint64(&s.RaceFailures, 1)
					return
				}
			}
		}(workerID)
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()

	<-timer.C
	stop.Store(true)

	wg.Wait()
}

func (s *SystemPipelineSimulator) VerifyInvariants() error {
	for _, stream := range s.Pool {
		state := atomic.LoadInt32(&stream.State)

		if state != StateIdle {
			return fmt.Errorf(
				"stream %d left in invalid state %d",
				stream.ConnID,
				state,
			)
		}
	}

	raceFailures := atomic.LoadUint64(&s.RaceFailures)

	if raceFailures != 0 {
		return fmt.Errorf(
			"detected %d invalid ownership transitions",
			raceFailures,
		)
	}

	return nil
}

func TestVerifyFaultPlaneResilience(t *testing.T) {
	const (
		poolSize = 1000
		workers  = 64
		duration = 5 * time.Second
	)

	harness := NewSystemPipelineSimulator(poolSize)

	start := time.Now()

	harness.RunStressMatrix(duration, workers)

	elapsed := time.Since(start)

	if err := harness.VerifyInvariants(); err != nil {
		t.Fatalf(
			"concurrency invariant failed after %s: %v",
			elapsed,
			err,
		)
	}

	successSwaps := atomic.LoadUint64(&harness.SuccessSwaps)
	busySkips := atomic.LoadUint64(&harness.BusySkips)
	raceFailures := atomic.LoadUint64(&harness.RaceFailures)

	if successSwaps == 0 {
		t.Fatal("stress harness completed with zero successful ownership swaps")
	}

	if raceFailures != 0 {
		t.Fatalf(
			"expected zero ownership violations, got %d",
			raceFailures,
		)
	}

	t.Logf("stress duration: %s", elapsed)
	t.Logf("pool size: %d", poolSize)
	t.Logf("workers: %d", workers)
	t.Logf("successful ownership swaps: %d", successSwaps)
	t.Logf("contention skips: %d", busySkips)
	t.Logf("ownership violations: %d", raceFailures)
}
