package runtime

import (
	"context"
	"log/slog"
	"testing"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

func TestLinuxKernelStateRuntimeRecoverValidation(t *testing.T) {
	runtime := &linuxKernelStateRuntime{}

	if err := runtime.Recover(context.Background(), "workflow-id"); err == nil {
		t.Fatal("expected error for uninitialized runtime")
	}

	runtime = &linuxKernelStateRuntime{
		enforcer: nil,
	}

	if err := runtime.Recover(context.Background(), "workflow-id"); err == nil {
		t.Fatal("expected error for nil enforcer")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runtime = &linuxKernelStateRuntime{
		enforcer: func() *kernel.ProductionKernelEnforcer {
			enforcer, err := kernel.NewProductionKernelEnforcer(
				slog.Default(),
				kernel.DefaultFaultPolicy{},
				kernel.NoopKernelAction{},
				&kernel.InMemoryStateTelemetry{},
			)
			if err != nil {
				t.Fatal(err)
			}
			return enforcer
		}(),
	}

	if err := runtime.Recover(ctx, "workflow-id"); err == nil {
		t.Fatal("expected context cancellation error")
	}

	if err := runtime.Recover(context.Background(), ""); err == nil {
		t.Fatal("expected empty workflow ID error")
	}
}

func TestLinuxKernelStateRuntimeRecoverAppliesRecoveryEvent(t *testing.T) {
	enforcer, err := kernel.NewProductionKernelEnforcer(
		slog.Default(),
		kernel.DefaultFaultPolicy{},
		kernel.NoopKernelAction{},
		&kernel.InMemoryStateTelemetry{},
	)
	if err != nil {
		t.Fatalf("create enforcer: %v", err)
	}

	if err := enforcer.Handle(context.Background(), kernel.FaultEvent{
		Type:   "connection_failure",
		Source: "test",
	}); err != nil {
		t.Fatalf("create degraded state: %v", err)
	}

	if got := enforcer.State(); got != kernel.StateDegraded {
		t.Fatalf("expected degraded state, got %s", got)
	}

	runtime := &linuxKernelStateRuntime{
		enforcer: enforcer,
		logger:   slog.Default(),
	}

	if err := runtime.Recover(context.Background(), "workflow-runtime-recovery"); err != nil {
		t.Fatalf("runtime recovery: %v", err)
	}

	if got := enforcer.State(); got != kernel.StateHealthy {
		t.Fatalf("expected healthy state after recovery, got %s", got)
	}
}
