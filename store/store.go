package store

import (
	"context"
	"time"
)

type StoreData interface {
	GetID() string
	GetExpiresAt() time.Time
}

type Store[T StoreData] interface {
	Set(ctx context.Context, session T) error
	Get(ctx context.Context, id string) (*T, error)
	GetByFilter(ctx context.Context, filter map[string]any, limit, offset int) ([]*T, error)
	Delete(ctx context.Context, id string) error
}
