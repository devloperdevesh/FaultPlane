package runtime

import (
	"context"
	"log/slog"

	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/gateway"
)

type Daemon struct {
	logger *slog.Logger

	control *control.Manager

	gateway *gateway.Manager
}

func New(
	logger *slog.Logger,
	control *control.Manager,
	gateway *gateway.Manager,
) *Daemon {

	return &Daemon{
		logger:  logger,
		control: control,
		gateway: gateway,
	}
}

func (d *Daemon) Start(ctx context.Context) error {

	d.logger.Info(
		"faultplane daemon starting",
	)

	go d.control.Start(ctx)

	go d.gateway.Start(ctx)

	<-ctx.Done()

	d.logger.Info(
		"faultplane daemon stopped",
	)

	return nil
}
