package wasm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tetratelabs/wazero"
)

var (
	ErrEmptyModule = errors.New("wasm module is empty")
	ErrTimeout     = errors.New("wasm execution timed out")
)

type Result struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   time.Duration
}

type Executor struct {
	runtime wazero.Runtime
}

func New(ctx context.Context) *Executor {
	return &Executor{
		runtime: wazero.NewRuntime(ctx),
	}
}

func (e *Executor) Close(ctx context.Context) error {
	if e == nil || e.runtime == nil {
		return nil
	}

	return e.runtime.Close(ctx)
}

func (e *Executor) Execute(
	ctx context.Context,
	module []byte,
	timeout time.Duration,
) (Result, error) {
	if len(module) == 0 {
		return Result{}, ErrEmptyModule
	}

	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	if timeout <= 0 {
		return Result{}, fmt.Errorf("invalid timeout: %w", ErrTimeout)
	}

	start := time.Now()

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	compiled, err := e.runtime.CompileModule(execCtx, module)
	if err != nil {
		if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
			return Result{}, ErrTimeout
		}

		return Result{}, fmt.Errorf("compile wasm module: %w", err)
	}

	defer compiled.Close(execCtx)

	_, err = e.runtime.InstantiateModule(

		execCtx,

		compiled,

		wazero.NewModuleConfig(),
	)
	if err != nil {
		if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
			return Result{}, ErrTimeout
		}

		return Result{}, fmt.Errorf("instantiate wasm module: %w", err)
	}

	finished := time.Now()

	return Result{
		StartedAt:  start,
		FinishedAt: finished,
		Duration:   finished.Sub(start),
	}, nil
}
