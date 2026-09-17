// SPDX-License-Identifier: Apache-2.0
//go:build linux

package kernel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// EnforcerState represents the lifecycle state of a kernel-controlled
// FaultPlane workload/resource.
type EnforcerState uint8

const (
	StateUnknown EnforcerState = iota
	StateHealthy
	StateDegraded
	StateRecovering
	StateEnforced
	StateFailed
)

func (s EnforcerState) String() string {
	switch s {
	case StateHealthy:
		return "healthy"
	case StateDegraded:
		return "degraded"
	case StateRecovering:
		return "recovering"
	case StateEnforced:
		return "enforced"
	case StateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// EnforcementAction is the concrete action selected by policy.
type EnforcementAction string

const (
	ActionObserve    EnforcementAction = "observe"
	ActionRecover    EnforcementAction = "recover"
	ActionEnforce    EnforcementAction = "enforce"
	ActionFailOpen   EnforcementAction = "fail_open"
	ActionFailClosed EnforcementAction = "fail_closed"
)

// StateTransition describes one accepted state transition.
type StateTransition struct {
	From      EnforcerState
	To        EnforcerState
	Action    EnforcementAction
	EventType string
	Reason    string
	At        time.Time
}

// StatePolicy evaluates an observed kernel event and decides the
// desired state/action.
type StatePolicy interface {
	Evaluate(FaultEvent, EnforcerState) (EnforcerState, EnforcementAction, error)
}

// StateEnforcerAction executes the enforcement/recovery operation.
// Concrete eBPF/kernel operations can be injected later without
// coupling the state machine to the mechanism.
type StateEnforcerAction interface {
	Execute(context.Context, EnforcementAction, FaultEvent) error
}

// StateTelemetry receives state-machine telemetry.
type StateTelemetry interface {
	RecordTransition(StateTransition)
	RecordError(error)
}

// ProductionKernelEnforcer is the synchronized state-machine boundary
// between kernel observations, policy decisions, enforcement actions,
// and telemetry.
//
// The state machine itself is mechanism-independent. Actual kernel/eBPF
// enforcement is injected through StateEnforcerAction.
type ProductionKernelEnforcer struct {
	mu sync.RWMutex

	logger *slog.Logger

	state EnforcerState

	policy    StatePolicy
	action    StateEnforcerAction
	telemetry StateTelemetry

	transitionCount uint64
	errorCount      uint64

	lastTransition StateTransition
}

// NewProductionKernelEnforcer creates a production state enforcer.
func NewProductionKernelEnforcer(
	logger *slog.Logger,
	policy StatePolicy,
	action StateEnforcerAction,
	telemetry StateTelemetry,
) (*ProductionKernelEnforcer, error) {
	if policy == nil {
		return nil, errors.New("state enforcer: policy is required")
	}

	if action == nil {
		return nil, errors.New("state enforcer: enforcement action is required")
	}

	if telemetry == nil {
		return nil, errors.New("state enforcer: telemetry is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &ProductionKernelEnforcer{
		logger:    logger,
		state:     StateHealthy,
		policy:    policy,
		action:    action,
		telemetry: telemetry,
	}, nil
}

// State returns the current state using a read lock.
func (e *ProductionKernelEnforcer) State() EnforcerState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.state
}

// Stats returns synchronized state-machine counters.
func (e *ProductionKernelEnforcer) Stats() (transitions, errors uint64) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.transitionCount, e.errorCount
}

// LastTransition returns a copy of the latest accepted transition.
func (e *ProductionKernelEnforcer) LastTransition() (StateTransition, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.transitionCount == 0 {
		return StateTransition{}, false
	}

	return e.lastTransition, true
}

// Handle observes one kernel FaultEvent, evaluates policy, performs
// the selected action, and commits the state transition only after
// the enforcement operation succeeds.
//
// This ordering is deliberate:
//
// observation
//
//	-> policy
//	-> enforcement
//	-> state commit
//	-> telemetry
//
// A failed enforcement therefore cannot silently appear as a
// successful state transition.
func (e *ProductionKernelEnforcer) Handle(
	ctx context.Context,
	event FaultEvent,
) error {
	if ctx == nil {
		return errors.New("state enforcer: nil context")
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("state enforcer: context unavailable: %w", err)
	}

	if event.Type == "" {
		return errors.New("state enforcer: event type is required")
	}

	// Snapshot state under lock. Policy evaluation is intentionally
	// performed outside the lock so slow policy implementations do
	// not block readers or unrelated state inspection.
	e.mu.RLock()
	current := e.state
	e.mu.RUnlock()

	next, action, err := e.policy.Evaluate(event, current)
	if err != nil {
		e.recordError(err)
		return fmt.Errorf("state enforcer: policy evaluation: %w", err)
	}

	if next == StateUnknown {
		err := errors.New("state enforcer: policy returned unknown state")
		e.recordError(err)
		return err
	}

	if !validTransition(current, next) {
		err := fmt.Errorf(
			"state enforcer: invalid transition %s -> %s",
			current,
			next,
		)
		e.recordError(err)
		return err
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("state enforcer: cancelled before enforcement: %w", err)
	}

	// Execute the real enforcement mechanism before committing state.
	if err := e.action.Execute(ctx, action, event); err != nil {
		e.recordError(err)
		return fmt.Errorf("state enforcer: enforcement action %q: %w", action, err)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("state enforcer: cancelled before state commit: %w", err)
	}

	transition := StateTransition{
		From:      current,
		To:        next,
		Action:    action,
		EventType: event.Type,
		Reason:    event.Reason,
		At:        time.Now(),
	}

	// Commit is serialized. Re-check the observed state so two
	// concurrent handlers cannot both overwrite an intervening state.
	e.mu.Lock()

	if e.state != current {
		actual := e.state
		e.mu.Unlock()

		err := fmt.Errorf(
			"state enforcer: stale transition %s -> %s; current state is %s",
			current,
			next,
			actual,
		)
		e.recordError(err)
		return err
	}

	e.state = next
	e.transitionCount++
	e.lastTransition = transition

	e.mu.Unlock()

	e.telemetry.RecordTransition(transition)

	e.logger.Debug(
		"kernel state transition committed",
		"from", transition.From.String(),
		"to", transition.To.String(),
		"action", transition.Action,
		"event_type", transition.EventType,
		"reason", transition.Reason,
	)

	return nil
}

func (e *ProductionKernelEnforcer) recordError(err error) {
	e.mu.Lock()
	e.errorCount++
	e.mu.Unlock()

	e.telemetry.RecordError(err)

	e.logger.Error(
		"kernel state enforcement error",
		"error", err,
	)
}

// validTransition defines the legal state-machine edges.
//
// Healthy:
//
//	healthy -> healthy/degraded
//
// Degraded:
//
//	degraded -> degraded/recovering/enforced/failed
//
// Recovering:
//
//	recovering -> recovering/healthy/enforced/failed
//
// Enforced:
//
//	enforced -> enforced/recovering/healthy/failed
//
// Failed:
//
//	failed -> failed/recovering
//
// Unknown is never a valid runtime state.
func validTransition(from, to EnforcerState) bool {
	switch from {
	case StateHealthy:
		return to == StateHealthy ||
			to == StateDegraded

	case StateDegraded:
		return to == StateDegraded ||
			to == StateRecovering ||
			to == StateEnforced ||
			to == StateHealthy ||
			to == StateFailed

	case StateRecovering:
		return to == StateRecovering ||
			to == StateHealthy ||
			to == StateEnforced ||
			to == StateFailed

	case StateEnforced:
		return to == StateEnforced ||
			to == StateRecovering ||
			to == StateHealthy ||
			to == StateFailed

	case StateFailed:
		return to == StateFailed ||
			to == StateRecovering

	default:
		return false
	}
}

// DefaultFaultPolicy provides a conservative production baseline.
// More sophisticated policy engines can implement StatePolicy later.
type DefaultFaultPolicy struct{}

func (DefaultFaultPolicy) Evaluate(
	event FaultEvent,
	current EnforcerState,
) (EnforcerState, EnforcementAction, error) {
	switch event.Type {
	case "connection_failure":
		switch current {
		case StateHealthy:
			return StateDegraded, ActionObserve, nil

		case StateDegraded, StateRecovering:
			return StateEnforced, ActionEnforce, nil

		case StateEnforced:
			return StateEnforced, ActionEnforce, nil

		case StateFailed:
			return StateRecovering, ActionRecover, nil
		}

	case "connection_recovered", "recovery_success":
		switch current {
		case StateDegraded, StateRecovering, StateEnforced:
			return StateHealthy, ActionRecover, nil

		case StateHealthy:
			return StateHealthy, ActionObserve, nil
		}

	case "kernel_error":
		return StateFailed, ActionFailClosed, nil

	case "kernel_warning", "connection_state_change":
		return current, ActionObserve, nil

	default:
		return current, ActionObserve, nil
	}

	return StateUnknown, ActionObserve,
		fmt.Errorf(
			"no policy for event %q in state %s",
			event.Type,
			current,
		)
}

// NoopKernelAction is a safe default mechanism for tests and
// development. Production kernel/eBPF implementations can replace it.
type NoopKernelAction struct{}

func (NoopKernelAction) Execute(
	ctx context.Context,
	_ EnforcementAction,
	_ FaultEvent,
) error {
	if ctx == nil {
		return errors.New("kernel action: nil context")
	}

	return ctx.Err()
}

// InMemoryStateTelemetry is a concurrency-safe telemetry sink useful
// for tests and local runtime validation.
type InMemoryStateTelemetry struct {
	mu sync.Mutex

	Transitions []StateTransition
	Errors      []error
}

func (t *InMemoryStateTelemetry) RecordTransition(
	transition StateTransition,
) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Transitions = append(t.Transitions, transition)
}

func (t *InMemoryStateTelemetry) RecordError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Errors = append(t.Errors, err)
}

func (t *InMemoryStateTelemetry) Snapshot() (
	[]StateTransition,
	[]error,
) {
	t.mu.Lock()
	defer t.mu.Unlock()

	transitions := append(
		[]StateTransition(nil),
		t.Transitions...,
	)

	errs := append(
		[]error(nil),
		t.Errors...,
	)

	return transitions, errs
}
