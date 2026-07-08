package storage

import (
	"context"
	"errors"
)

var ErrStateNotFound = errors.New("state signature not found")

type Store interface {
	Save(ctx context.Context, agentID string, step uint64, payload string, anchor string) error
	Get(ctx context.Context, agentID string) (string, uint64, string, error)
}
