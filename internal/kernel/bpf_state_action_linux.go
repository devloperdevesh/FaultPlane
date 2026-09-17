//go:build linux

package kernel

import (
	"context"
	"fmt"
	"log/slog"
)

// BPFKernelAction connects the state enforcer to the live
// eBPF enforcement control map.
type BPFKernelAction struct {
	logger *slog.Logger
	loader *Loader
}

func NewBPFKernelAction(
	logger *slog.Logger,
	loader *Loader,
) (*BPFKernelAction, error) {
	if loader == nil {
		return nil, fmt.Errorf("BPF kernel action: loader is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &BPFKernelAction{
		logger: logger,
		loader: loader,
	}, nil
}

func (a *BPFKernelAction) Execute(
	ctx context.Context,
	action EnforcementAction,
	event FaultEvent,
) error {
	if ctx == nil {
		return fmt.Errorf("BPF kernel action: context is nil")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	var mode uint32

	switch action {
	case ActionObserve:
		mode = EnforcementAllow

	case ActionRecover:
		mode = EnforcementAllow

	case ActionEnforce:
		mode = EnforcementRedirect

	case ActionFailOpen:
		mode = EnforcementAllow

	case ActionFailClosed:
		mode = EnforcementDrop

	default:
		return fmt.Errorf("unsupported enforcement action: %q", action)
	}

	if err := a.loader.SetEnforcementMode(mode); err != nil {
		return fmt.Errorf(
			"apply enforcement action %q for event %q: %w",
			action,
			event.Type,
			err,
		)
	}

	a.logger.Debug(
		"eBPF enforcement action applied",
		"action", action,
		"mode", mode,
		"event_type", event.Type,
		"source", event.Source,
		"reason", event.Reason,
	)

	return nil
}
