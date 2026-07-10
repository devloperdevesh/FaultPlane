package storage

import (
	"context"

	"github.com/devloperdevesh/agentmesh/internal/control"
)

type Store interface {
	Save(ctx context.Context, workflow *control.Workflow) error

	Load(ctx context.Context, id string) (*control.Workflow, error)

	Delete(ctx context.Context, id string) error

	List(ctx context.Context) ([]*control.Workflow, error)
}