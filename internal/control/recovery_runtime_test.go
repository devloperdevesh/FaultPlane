package control

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/devloperdevesh/FaultPlane/internal/domain"
	"github.com/devloperdevesh/FaultPlane/internal/storage"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type recordingRuntimeExecutor struct {
	mu         sync.Mutex
	calls      []string
	recoverErr error
}

func (r *recordingRuntimeExecutor) Recover(ctx context.Context, workflowID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.calls = append(r.calls, workflowID)

	if ctx == nil {
		return errors.New("context is nil")
	}

	return r.recoverErr
}

func (r *recordingRuntimeExecutor) Calls() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]string(nil), r.calls...)
}

func newRuntimeRecoveryController(runtime RuntimeExecutor) (*Controller, string, error) {
	store := storage.NewMemoryStore()
	collector := telemetry.NewCollector(telemetry.NewRegistry())
	controller := NewController(store, collector, runtime)

	workflow := &domain.Workflow{
		ID:   "workflow-runtime-integration",
		Name: "runtime integration",
	}

	if err := controller.Register(workflow); err != nil {
		return nil, "", err
	}

	if err := controller.CreateCheckpoint(workflow.ID, 1, []byte("checkpoint")); err != nil {
		return nil, "", err
	}

	return controller, workflow.ID, nil
}

func TestRecoverContextExecutesRuntimeBeforeResume(t *testing.T) {
	runtime := &recordingRuntimeExecutor{}

	controller, workflowID, err := newRuntimeRecoveryController(runtime)
	if err != nil {
		t.Fatalf("setup controller: %v", err)
	}

	if err := controller.RecoverContext(context.Background(), workflowID); err != nil {
		t.Fatalf("recover context: %v", err)
	}

	calls := runtime.Calls()

	if len(calls) != 1 {
		t.Fatalf("expected exactly one runtime recovery call, got %d", len(calls))
	}

	if calls[0] != workflowID {
		t.Fatalf("expected runtime recovery for %q, got %q", workflowID, calls[0])
	}

	recovered, err := controller.Get(workflowID)
	if err != nil {
		t.Fatalf("get recovered workflow: %v", err)
	}

	if recovered.Status != StatusRunning {
		t.Fatalf(
			"expected workflow status %q after recovery, got %q",
			StatusRunning,
			recovered.Status,
		)
	}
}

func TestRecoverContextStopsBeforeResumeWhenRuntimeFails(t *testing.T) {
	runtime := &recordingRuntimeExecutor{
		recoverErr: errors.New("runtime recovery failed"),
	}

	controller, workflowID, err := newRuntimeRecoveryController(runtime)
	if err != nil {
		t.Fatalf("setup controller: %v", err)
	}

	err = controller.RecoverContext(context.Background(), workflowID)
	if err == nil {
		t.Fatal("expected runtime recovery error")
	}

	calls := runtime.Calls()

	if len(calls) != 1 {
		t.Fatalf("expected exactly one runtime recovery call, got %d", len(calls))
	}

	recovered, err := controller.Get(workflowID)
	if err != nil {
		t.Fatalf("get workflow: %v", err)
	}

	if recovered.Status == StatusRunning {
		t.Fatal("workflow must not resume when runtime recovery fails")
	}
}
