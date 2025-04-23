package authz

import "context"

type Repository[T any, K comparable] interface {
	GetByID(ctx context.Context, id K) (T, error)
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id K) error
}
