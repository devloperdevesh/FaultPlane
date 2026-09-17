// SPDX-License-Identifier: Apache-2.0
//go:build linux

package kernel

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
)

type testKernelAction struct {
	mu       sync.Mutex
	actions  []EnforcementAction
	failWith error
}

func (a *testKernelAction) Execute(
	ctx context.Context,
	action EnforcementAction,
	_ FaultEvent,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.actions = append(a.actions, action)

	return a.failWith
}

func TestProductionKernelEnforcer_InitialState(t *testing.T) {
	telemetry := &InMemoryStateTelemetry{}
	action := &testKernelAction{}

	enforcer, err := NewProductionKernelEnforcer(
		slog.Default(),
		DefaultFaultPolicy{},
		action,
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	if got := enforcer.State(); got != StateHealthy {
		t.Fatalf("expected initial healthy state, got %s", got)
	}

	transitions, errs := enforcer.Stats()

	if transitions != 0 {
		t.Fatalf("expected zero transitions, got %d", transitions)
	}

	if errs != 0 {
		t.Fatalf("expected zero errors, got %d", errs)
	}
}

func TestProductionKernelEnforcer_TransitionFlow(t *testing.T) {
	telemetry := &InMemoryStateTelemetry{}
	action := &testKernelAction{}

	enforcer, err := NewProductionKernelEnforcer(
		slog.Default(),
		DefaultFaultPolicy{},
		action,
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	event := FaultEvent{
		Type:   "connection_failure",
		Source: "kernel/tcp",
		Reason: "CLOSED",
	}

	if err := enforcer.Handle(context.Background(), event); err != nil {
		t.Fatal(err)
	}

	if got := enforcer.State(); got != StateDegraded {
		t.Fatalf("expected degraded state, got %s", got)
	}

	transitions, errs := enforcer.Stats()

	if transitions != 1 {
		t.Fatalf("expected one transition, got %d", transitions)
	}

	if errs != 0 {
		t.Fatalf("expected zero errors, got %d", errs)
	}

	recorded, telemetryErrors := telemetry.Snapshot()

	if len(recorded) != 1 {
		t.Fatalf("expected one telemetry transition, got %d", len(recorded))
	}

	if len(telemetryErrors) != 0 {
		t.Fatalf("expected zero telemetry errors, got %d", len(telemetryErrors))
	}

	if recorded[0].From != StateHealthy ||
		recorded[0].To != StateDegraded {
		t.Fatalf(
			"unexpected transition: %s -> %s",
			recorded[0].From,
			recorded[0].To,
		)
	}
}

func TestProductionKernelEnforcer_EnforcementFailureDoesNotCommit(
	t *testing.T,
) {
	telemetry := &InMemoryStateTelemetry{}
	action := &testKernelAction{
		failWith: errors.New("simulated kernel enforcement failure"),
	}

	enforcer, err := NewProductionKernelEnforcer(
		slog.Default(),
		DefaultFaultPolicy{},
		action,
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = enforcer.Handle(
		context.Background(),
		FaultEvent{
			Type:   "connection_failure",
			Source: "kernel/tcp",
			Reason: "CLOSED",
		},
	)

	if err == nil {
		t.Fatal("expected enforcement error")
	}

	if got := enforcer.State(); got != StateHealthy {
		t.Fatalf(
			"state committed despite enforcement failure: %s",
			got,
		)
	}

	transitions, errorsCount := enforcer.Stats()

	if transitions != 0 {
		t.Fatalf(
			"expected zero committed transitions, got %d",
			transitions,
		)
	}

	if errorsCount != 1 {
		t.Fatalf(
			"expected one recorded error, got %d",
			errorsCount,
		)
	}
}

func TestProductionKernelEnforcer_ContextCancellation(
	t *testing.T,
) {
	telemetry := &InMemoryStateTelemetry{}
	action := &testKernelAction{}

	enforcer, err := NewProductionKernelEnforcer(
		slog.Default(),
		DefaultFaultPolicy{},
		action,
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = enforcer.Handle(
		ctx,
		FaultEvent{
			Type:   "connection_failure",
			Source: "kernel/tcp",
		},
	)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if got := enforcer.State(); got != StateHealthy {
		t.Fatalf("expected healthy state, got %s", got)
	}
}

func TestProductionKernelEnforcer_InvalidEvent(t *testing.T) {
	telemetry := &InMemoryStateTelemetry{}
	action := &testKernelAction{}

	enforcer, err := NewProductionKernelEnforcer(
		slog.Default(),
		DefaultFaultPolicy{},
		action,
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = enforcer.Handle(
		context.Background(),
		FaultEvent{},
	)

	if err == nil {
		t.Fatal("expected invalid event error")
	}

	if got := enforcer.State(); got != StateHealthy {
		t.Fatalf("expected healthy state, got %s", got)
	}
}

func TestProductionKernelEnforcer_ConcurrentAccess(t *testing.T) {
	telemetry := &InMemoryStateTelemetry{}
	action := &testKernelAction{}

	enforcer, err := NewProductionKernelEnforcer(
		slog.Default(),
		DefaultFaultPolicy{},
		action,
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	const workers = 32
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(workers)

	for worker := 0; worker < workers; worker++ {
		go func() {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				_ = enforcer.State()

				_, _ = enforcer.Stats()

				_, _ = enforcer.LastTransition()
			}
		}()
	}

	wg.Wait()

	if got := enforcer.State(); got == StateUnknown {
		t.Fatal("enforcer returned unknown runtime state")
	}
}

func TestValidTransitionMatrix(t *testing.T) {
	valid := []struct {
		from EnforcerState
		to   EnforcerState
	}{
		{StateHealthy, StateHealthy},
		{StateHealthy, StateDegraded},
		{StateDegraded, StateRecovering},
		{StateDegraded, StateEnforced},
		{StateDegraded, StateFailed},
		{StateRecovering, StateHealthy},
		{StateRecovering, StateEnforced},
		{StateRecovering, StateFailed},
		{StateEnforced, StateRecovering},
		{StateEnforced, StateHealthy},
		{StateEnforced, StateFailed},
		{StateFailed, StateRecovering},
	}

	for _, tc := range valid {
		if !validTransition(tc.from, tc.to) {
			t.Fatalf(
				"expected valid transition %s -> %s",
				tc.from,
				tc.to,
			)
		}
	}

	invalid := []struct {
		from EnforcerState
		to   EnforcerState
	}{
		{StateUnknown, StateHealthy},
		{StateHealthy, StateRecovering},
		{StateHealthy, StateEnforced},
		{StateHealthy, StateFailed},
		{StateFailed, StateHealthy},
		{StateFailed, StateEnforced},
	}

	for _, tc := range invalid {
		if validTransition(tc.from, tc.to) {
			t.Fatalf(
				"expected invalid transition %s -> %s",
				tc.from,
				tc.to,
			)
		}
	}
}
