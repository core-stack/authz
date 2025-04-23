package zrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/core-stack/authz/zmodel"
)

type IOauth2UserRepository interface {
	GetByID(ctx context.Context, id string) (*zmodel.OAuth2User, error)
	FindByProviderID(ctx context.Context, provider string, providerID string) (*zmodel.OAuth2User, error)
	CreateWithUser(ctx context.Context, o *zmodel.OAuth2User) (*zmodel.OAuth2User, error)
	Create(ctx context.Context, o *zmodel.OAuth2User) error
	Update(ctx context.Context, o *zmodel.OAuth2User) error
	Delete(ctx context.Context, id string) error
	Transaction(ctx context.Context, fn func(context.Context) error) error
}

type OAuth2UserRepository struct {
	db *sql.DB
}

func NewOAuth2UserRepository(db *sql.DB) *OAuth2UserRepository {
	return &OAuth2UserRepository{db: db}
}

func (r *OAuth2UserRepository) GetByID(ctx context.Context, id string) (*zmodel.OAuth2User, error) {
	return r.queryOne(ctx, "SELECT id, user_id, provider, provider_id, linked_at FROM oauth2_users WHERE id = $1", id)
}

func (r *OAuth2UserRepository) FindByProviderID(ctx context.Context, provider string, providerID string) (*zmodel.OAuth2User, error) {
	return r.queryOne(ctx, "SELECT id, user_id, provider, provider_id, linked_at FROM oauth2_users WHERE provider = $1 AND provider_id = $2", provider, providerID)
}

func (r *OAuth2UserRepository) Create(ctx context.Context, o *zmodel.OAuth2User) error {
	_, err := exec(ctx, r.db, `
		INSERT INTO oauth2_users (id, user_id, provider, provider_id, linked_at)
		VALUES ($1, $2, $3, $4, $5)
	`, o.ID, o.UserID, o.Provider, o.ProviderID, o.LinkedAt)
	return err
}

func (r *OAuth2UserRepository) CreateWithUser(ctx context.Context, o *zmodel.OAuth2User) (*zmodel.OAuth2User, error) {
	err := r.Transaction(ctx, func(txCtx context.Context) error {
		// cria o usuário se necessário
		if o.User != nil {
			err := execUserInsert(txCtx, r.db, o.User)
			if err != nil {
				return err
			}
			o.UserID = o.User.ID
		}
		return r.Create(txCtx, o)
	})
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *OAuth2UserRepository) Update(ctx context.Context, o *zmodel.OAuth2User) error {
	_, err := exec(ctx, r.db, `
		UPDATE oauth2_users
		SET user_id = $2, provider = $3, provider_id = $4, linked_at = $5
		WHERE id = $1
	`, o.ID, o.UserID, o.Provider, o.ProviderID, o.LinkedAt)
	return err
}

func (r *OAuth2UserRepository) Delete(ctx context.Context, id string) error {
	_, err := exec(ctx, r.db, "DELETE FROM oauth2_users WHERE id = $1", id)
	return err
}

func (r *OAuth2UserRepository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}

// --- helpers ---

func (r *OAuth2UserRepository) queryOne(ctx context.Context, query string, args ...any) (*zmodel.OAuth2User, error) {
	row := queryRow(ctx, r.db, query, args...)
	var o zmodel.OAuth2User
	err := row.Scan(&o.ID, &o.UserID, &o.Provider, &o.ProviderID, &o.LinkedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

// --- optional helper for creating user ---

func execUserInsert(ctx context.Context, db *sql.DB, user *zmodel.User) error {
	_, err := exec(ctx, db, `
		INSERT INTO users (id, email, name, created_at)
		VALUES ($1, $2, $3, $4)
	`, user.ID, user.Email, user.Name, user.CreatedAt)
	return err
}
