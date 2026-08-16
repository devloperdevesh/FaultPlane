package kernel

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNetlinkListener_StartAndShutdown(t *testing.T) {

	logger := slog.Default()

	listener := NewNetlinkListener(
		logger,
	)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	err := listener.Start(ctx)

	if err != nil {
		t.Fatalf(
			"start failed: %v",
			err,
		)
	}

	cancel()

	time.Sleep(
		50 * time.Millisecond,
	)

}

func TestNetlinkListener_DuplicateStart(t *testing.T) {

	logger := slog.Default()

	listener := NewNetlinkListener(
		logger,
	)

	ctx := context.Background()

	err := listener.Start(ctx)

	if err != nil {
		t.Fatal(err)
	}

	err = listener.Start(ctx)

	if err == nil {
		t.Fatal(
			"expected duplicate start error",
		)
	}

}
