package zrepository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/core-stack/authz/zmodel"
	"github.com/core-stack/authz/zrepository"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*zrepository.CodeRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := zrepository.NewCodeRepository(db)
	return repo, mock, func() { db.Close() }
}

func TestCodeRepository_GetByID(t *testing.T) {
	repo, mock, teardown := setup(t)
	defer teardown()

	now := time.Now()
	expected := &zmodel.Code{
		ID:        "123",
		UserID:    "u1",
		Token:     "token",
		Type:      zmodel.ActiveAccount,
		ExpiresAt: now,
		CreatedAt: now,
		UsedAt:    sql.NullTime{},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "token", "type", "expires_at", "created_at", "used_at"}).
		AddRow(expected.ID, expected.UserID, expected.Token, expected.Type, expected.ExpiresAt, expected.CreatedAt, expected.UsedAt)

	mock.ExpectQuery("SELECT id, user_id, token, type, expires_at, created_at, used_at FROM codes WHERE id = \\$1").
		WithArgs("123").WillReturnRows(rows)

	actual, err := repo.GetByID(context.Background(), "123")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
func TestCodeRepository_GetByToken(t *testing.T) {
	repo, mock, teardown := setup(t)
	defer teardown()

	now := time.Now()
	expected := &zmodel.Code{
		ID:        "999",
		UserID:    "u42",
		Token:     "magic-token",
		Type:      zmodel.ActiveAccount,
		ExpiresAt: now,
		CreatedAt: now,
		UsedAt:    sql.NullTime{},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "token", "type", "expires_at", "created_at", "used_at"}).
		AddRow(expected.ID, expected.UserID, expected.Token, expected.Type, expected.ExpiresAt, expected.CreatedAt, expected.UsedAt)

	mock.ExpectQuery("SELECT id, user_id, token, type, expires_at, created_at, used_at FROM codes WHERE token = \\$1").
		WithArgs("magic-token").WillReturnRows(rows)

	actual, err := repo.GetByToken(context.Background(), "magic-token")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
func TestCodeRepository_Create(t *testing.T) {
	repo, mock, teardown := setup(t)
	defer teardown()

	code := &zmodel.Code{
		ID:        "456",
		UserID:    "u2",
		Token:     "t2",
		Type:      zmodel.ActiveAccount,
		ExpiresAt: time.Now(),
		CreatedAt: time.Now(),
		UsedAt:    sql.NullTime{},
	}

	mock.ExpectExec("INSERT INTO codes").
		WithArgs(code.ID, code.UserID, code.Token, code.Type, code.ExpiresAt, code.CreatedAt, code.UsedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), code)
	require.NoError(t, err)
}

func TestCodeRepository_Update(t *testing.T) {
	repo, mock, teardown := setup(t)
	defer teardown()

	code := &zmodel.Code{
		ID:        "xyz",
		UserID:    "u88",
		Token:     "new-token",
		Type:      zmodel.ActiveAccount,
		ExpiresAt: time.Now(),
		CreatedAt: time.Now(),
		UsedAt:    sql.NullTime{},
	}

	mock.ExpectExec("UPDATE codes").
		WithArgs(code.ID, code.UserID, code.Token, code.Type, code.ExpiresAt, code.CreatedAt, code.UsedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), code)
	require.NoError(t, err)
}

func TestCodeRepository_Delete(t *testing.T) {
	repo, mock, teardown := setup(t)
	defer teardown()

	mock.ExpectExec("DELETE FROM codes WHERE id = \\$1").
		WithArgs("123").WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), "123")
	require.NoError(t, err)
}

func TestCodeRepository_Transaction(t *testing.T) {
	repo, mock, teardown := setup(t)
	defer teardown()

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := repo.Transaction(context.Background(), func(ctx context.Context) error {
		// Você poderia fazer algo aqui com exec/insert/update, mockado também.
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
