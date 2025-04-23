package zrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/core-stack/authz/zmodel"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id string) (*zmodel.User, error)
	GetByEmail(ctx context.Context, email string) (*zmodel.User, error)
	GetByUsername(ctx context.Context, username string) (*zmodel.User, error)
	Create(ctx context.Context, u *zmodel.User) error
	Update(ctx context.Context, u *zmodel.User) error
	UpdatePassword(ctx context.Context, userID, newPassword string) error
	Delete(ctx context.Context, id string) error
	Transaction(ctx context.Context, fn func(context.Context) error) error
}
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*zmodel.User, error) {
	var u zmodel.User
	err := queryRow(ctx, r.db, `
		SELECT id, username, name, password, email, status, is_admin, role_id, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Username, &u.Name, &u.Password, &u.Email,
		&u.Status, &u.IsAdmin, &u.RoleID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*zmodel.User, error) {
	var u zmodel.User
	err := queryRow(ctx, r.db, `
		SELECT id, username, name, password, email, status, is_admin, role_id, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &u.Username, &u.Name, &u.Password, &u.Email,
		&u.Status, &u.IsAdmin, &u.RoleID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*zmodel.User, error) {
	var u zmodel.User
	err := queryRow(ctx, r.db, `
		SELECT id, username, name, password, email, status, is_admin, role_id, created_at, updated_at
		FROM users WHERE username = $1
	`, username).Scan(
		&u.ID, &u.Username, &u.Name, &u.Password, &u.Email,
		&u.Status, &u.IsAdmin, &u.RoleID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *zmodel.User) error {
	now := time.Now()
	query := `
		INSERT INTO users (username, name, password, email, status, is_admin, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		RETURNING id, created_at, updated_at
	`
	err := queryRow(ctx, r.db, query,
		u.Username,
		u.Name,
		u.Password,
		u.Email,
		u.Status,
		u.IsAdmin,
		u.RoleID,
		now,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *zmodel.User) error {
	u.UpdatedAt = time.Now()
	_, err := exec(ctx, r.db, `
		UPDATE users SET name=$1, username=$2, password=$3, email=$4, status=$5, is_admin=$6, updated_at=$7
		WHERE id = $8
	`, u.Name, u.Username, u.Password, u.Email, u.Status, u.IsAdmin, u.UpdatedAt, u.ID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	_, err := exec(ctx, r.db, `
		UPDATE users SET password = $1, updated_at = $2
		WHERE id = $3
	`, newPassword, time.Now(), userID)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := exec(ctx, r.db, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}
