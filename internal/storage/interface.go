package storage

import (
	"context"

	"github.com/devloperdevesh/agentmesh/internal/domain"
)

type Store interface {
	Save(
		ctx context.Context,
		workflow *domain.Workflow,
	) error

	Load(
		ctx context.Context,
		id string,
	) (*domain.Workflow, error)

	Delete(
		ctx context.Context,
		id string,
	) error

	List(
		ctx context.Context,
	) ([]*domain.Workflow, error)
}
