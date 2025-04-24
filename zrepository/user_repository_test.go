package zrepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/core-stack/authz/zmodel"
	"github.com/core-stack/authz/zrepository"
	"github.com/stretchr/testify/require"
)

func newTestUser() *zmodel.User {
	now := time.Now()
	return &zmodel.User{
		ID:        "123",
		Username:  "johndoe",
		Name:      "John Doe",
		Password:  "hashed-pass",
		Email:     "john@example.com",
		Status:    zmodel.Active,
		IsAdmin:   false,
		RoleID:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)
	user := newTestUser()

	mock.ExpectQuery("SELECT id, username, name, password, email").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "name", "password", "email", "status", "is_admin", "role_id", "created_at", "updated_at",
		}).AddRow(
			user.ID, user.Username, user.Name, user.Password, user.Email,
			user.Status, user.IsAdmin, user.RoleID, user.CreatedAt, user.UpdatedAt,
		))

	got, err := repo.GetByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)
	user := newTestUser()

	mock.ExpectQuery("SELECT id, username, name, password, email").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "name", "password", "email", "status", "is_admin", "role_id", "created_at", "updated_at",
		}).AddRow(
			user.ID, user.Username, user.Name, user.Password, user.Email,
			user.Status, user.IsAdmin, user.RoleID, user.CreatedAt, user.UpdatedAt,
		))

	got, err := repo.GetByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, got.Email)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)
	user := newTestUser()

	mock.ExpectQuery("SELECT id, username, name, password, email").
		WithArgs(user.Username).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "name", "password", "email", "status", "is_admin", "role_id", "created_at", "updated_at",
		}).AddRow(
			user.ID, user.Username, user.Name, user.Password, user.Email,
			user.Status, user.IsAdmin, user.RoleID, user.CreatedAt, user.UpdatedAt,
		))

	got, err := repo.GetByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, got.Username)
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)
	user := &zmodel.User{
		Username: "jane",
		Name:     "Jane Smith",
		Password: "hashedpass",
		Email:    "jane@example.com",
		Status:   zmodel.Active,
		IsAdmin:  true,
		RoleID:   1,
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Name, user.Password, user.Email, user.Status, user.IsAdmin, user.RoleID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("new-id", time.Now(), time.Now()))

	err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.Equal(t, "new-id", user.ID)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)
	user := newTestUser()
	user.Password = "updated-hash"

	mock.ExpectExec("UPDATE users SET name=\\$1").
		WithArgs(user.Name, user.Username, user.Password, user.Email, user.Status, user.IsAdmin, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), user)
	require.NoError(t, err)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)

	mock.ExpectExec("UPDATE users SET password").
		WithArgs("new-password", sqlmock.AnyArg(), "user-id").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdatePassword(context.Background(), "user-id", "new-password")
	require.NoError(t, err)
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)

	mock.ExpectExec("DELETE FROM users").
		WithArgs("user-id").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(context.Background(), "user-id")
	require.NoError(t, err)
}

func TestUserRepository_Transaction(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := zrepository.NewUserRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE users SET password").
		WithArgs("tx-password", sqlmock.AnyArg(), "tx-user").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Transaction(context.Background(), func(ctx context.Context) error {
		return repo.UpdatePassword(ctx, "tx-user", "tx-password")
	})
	require.NoError(t, err)
}
