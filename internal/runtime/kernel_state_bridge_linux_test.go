// SPDX-License-Identifier: Apache-2.0
//go:build linux

package runtime

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

func TestFaultEventFromNetlink(t *testing.T) {
	event := kernel.Event{
		Type: "link",
		Data: []byte{1, 2, 3, 4},
	}

	got := faultEventFromNetlink(event)

	if got.Type != "kernel_warning" {
		t.Fatalf("expected kernel_warning, got %q", got.Type)
	}

	if got.Source != "kernel/rtnetlink" {
		t.Fatalf("unexpected source: %q", got.Source)
	}

	if got.Reason != "link" {
		t.Fatalf("unexpected reason: %q", got.Reason)
	}

	if got.Metadata["source"] != "rtnetlink" {
		t.Fatalf("unexpected source metadata: %q", got.Metadata["source"])
	}

	if got.Metadata["event_type"] != "link" {
		t.Fatalf("unexpected event type metadata: %q", got.Metadata["event_type"])
	}

	if got.Metadata["payload_bytes"] != "4" {
		t.Fatalf("unexpected payload size: %q", got.Metadata["payload_bytes"])
	}
}

func TestKernelStateBridge_StartCancellation(t *testing.T) {
	listener := kernel.NewNetlinkListener(slog.Default())

	telemetry := &kernel.InMemoryStateTelemetry{}

	enforcer, err := kernel.NewProductionKernelEnforcer(
		slog.Default(),
		kernel.DefaultFaultPolicy{},
		kernel.NoopKernelAction{},
		telemetry,
	)
	if err != nil {
		t.Fatal(err)
	}

	bridge, err := NewKernelStateBridge(
		slog.Default(),
		listener,
		enforcer,
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)

	go func() {
		done <- bridge.Start(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("bridge returned error: %v", err)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("bridge did not stop after context cancellation")
	}
}
