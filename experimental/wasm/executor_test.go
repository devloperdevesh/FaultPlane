package wasm

import (
	"context"
	"testing"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func minimalWASM(t *testing.T) []byte {
	t.Helper()

	ctx := context.Background()
	runtime := wazero.NewRuntime(ctx)
	defer runtime.Close(ctx)

	module, err := runtime.CompileModule(ctx, []byte{
		0x00, 0x61, 0x73, 0x6d,
		0x01, 0x00, 0x00, 0x00,
	})
	if err != nil {
		t.Fatalf("compile minimal wasm: %v", err)
	}

	defer module.Close(ctx)

	return []byte{
		0x00, 0x61, 0x73, 0x6d,
		0x01, 0x00, 0x00, 0x00,
	}
}

func TestExecutorRejectsEmptyModule(t *testing.T) {
	executor := New(context.Background())
	defer executor.Close(context.Background())

	_, err := executor.Execute(
		context.Background(),
		nil,
		time.Second,
	)

	if err != ErrEmptyModule {
		t.Fatalf("expected ErrEmptyModule, got %v", err)
	}
}

func TestExecutorExecutesValidModule(t *testing.T) {
	executor := New(context.Background())
	defer executor.Close(context.Background())

	module := minimalWASM(t)

	result, err := executor.Execute(
		context.Background(),
		module,
		time.Second,
	)

	if err != nil {
		t.Fatalf("execute wasm: %v", err)
	}

	if result.Duration < 0 {
		t.Fatalf("invalid execution duration: %v", result.Duration)
	}

	if result.FinishedAt.Before(result.StartedAt) {
		t.Fatalf("finished timestamp precedes start timestamp")
	}
}

func TestExecutorRespectsCancellation(t *testing.T) {
	executor := New(context.Background())
	defer executor.Close(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := executor.Execute(
		ctx,
		[]byte{
			0x00, 0x61, 0x73, 0x6d,
			0x01, 0x00, 0x00, 0x00,
		},
		time.Second,
	)

	if err == nil {
		t.Fatal("expected cancellation error")
	}

	_ = api.Module(nil)
}
