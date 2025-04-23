package zrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/core-stack/authz/zmodel"
)

type ICodeRepository interface {
	GetByID(ctx context.Context, id string) (*zmodel.Code, error)
	GetByToken(ctx context.Context, token string) (*zmodel.Code, error)
	Create(ctx context.Context, u *zmodel.Code) error
	Update(ctx context.Context, u *zmodel.Code) error
	Delete(ctx context.Context, id string) error
	Transaction(ctx context.Context, fn func(context.Context) error) error
}

type CodeRepository struct {
	db *sql.DB
}

func NewCodeRepository(db *sql.DB) *CodeRepository {
	return &CodeRepository{db: db}
}

func (r *CodeRepository) GetByID(ctx context.Context, id string) (*zmodel.Code, error) {
	return r.queryOne(ctx, "SELECT id, user_id, token, type, expires_at, created_at, used_at FROM codes WHERE id = $1", id)
}

func (r *CodeRepository) GetByToken(ctx context.Context, token string) (*zmodel.Code, error) {
	return r.queryOne(ctx, "SELECT id, user_id, token, type, expires_at, created_at, used_at FROM codes WHERE token = $1", token)
}

func (r *CodeRepository) Create(ctx context.Context, c *zmodel.Code) error {
	_, err := exec(ctx, r.db, `
		INSERT INTO codes (id, user_id, token, type, expires_at, created_at, used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.ID, c.UserID, c.Token, c.Type, c.ExpiresAt, c.CreatedAt, c.UsedAt,
	)
	return err
}

func (r *CodeRepository) Update(ctx context.Context, c *zmodel.Code) error {
	_, err := exec(ctx, r.db, `
		UPDATE codes
		SET user_id = $2, token = $3, type = $4, expires_at = $5, created_at = $6, used_at = $7
		WHERE id = $1`,
		c.ID, c.UserID, c.Token, c.Type, c.ExpiresAt, c.CreatedAt, c.UsedAt,
	)
	return err
}

func (r *CodeRepository) Delete(ctx context.Context, id string) error {
	_, err := exec(ctx, r.db, "DELETE FROM codes WHERE id = $1", id)
	return err
}

func (r *CodeRepository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}

func (r *CodeRepository) queryOne(ctx context.Context, query string, args ...any) (*zmodel.Code, error) {
	row := queryRow(ctx, r.db, query, args...)
	var c zmodel.Code
	err := row.Scan(&c.ID, &c.UserID, &c.Token, &c.Type, &c.ExpiresAt, &c.CreatedAt, &c.UsedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
