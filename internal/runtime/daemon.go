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

	go func() {

		err := d.control.Start(ctx)

		if err != nil {
			d.logger.Error(
				"control manager failed",
				"error",
				err,
			)
		}

	}()

	go func() {

		err := d.gateway.Start(ctx)

		if err != nil {
			d.logger.Error(
				"gateway manager failed",
				"error",
				err,
			)
		}

	}()

	<-ctx.Done()

	d.logger.Info(
		"faultplane daemon stopped",
	)

	return nil
}
