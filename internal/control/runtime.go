package control

import "context"

// RuntimeExecutor is the narrow control-plane boundary into the runtime.
// The control package owns the interface; runtime provides the implementation.
type RuntimeExecutor interface {
	Recover(context.Context, string) error
}
