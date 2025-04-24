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

func setupOAuth2(t *testing.T) (*zrepository.OAuth2UserRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return zrepository.NewOAuth2UserRepository(db), mock, func() {
		db.Close()
	}
}

func TestOAuth2UserRepository_GetByID(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	now := time.Now()
	expected := &zmodel.OAuth2User{
		ID:         "123",
		UserID:     "u1",
		Provider:   "google",
		ProviderID: "g123",
		LinkedAt:   now,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "provider", "provider_id", "linked_at"}).
		AddRow(expected.ID, expected.UserID, expected.Provider, expected.ProviderID, expected.LinkedAt)

	mock.ExpectQuery("SELECT id, user_id, provider, provider_id, linked_at FROM oauth2_users WHERE id = \\$1").
		WithArgs("123").WillReturnRows(rows)

	got, err := repo.GetByID(context.Background(), "123")
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestOAuth2UserRepository_FindByProviderID(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	now := time.Now()
	expected := &zmodel.OAuth2User{
		ID:         "123",
		UserID:     "u1",
		Provider:   "github",
		ProviderID: "gh_42",
		LinkedAt:   now,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "provider", "provider_id", "linked_at"}).
		AddRow(expected.ID, expected.UserID, expected.Provider, expected.ProviderID, expected.LinkedAt)

	mock.ExpectQuery("SELECT id, user_id, provider, provider_id, linked_at FROM oauth2_users WHERE provider = \\$1 AND provider_id = \\$2").
		WithArgs("github", "gh_42").WillReturnRows(rows)

	got, err := repo.FindByProviderID(context.Background(), "github", "gh_42")
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestOAuth2UserRepository_Create(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	o := &zmodel.OAuth2User{
		ID:         "abc",
		UserID:     "u42",
		Provider:   "github",
		ProviderID: "gh_abc",
		LinkedAt:   time.Now(),
	}

	mock.ExpectExec("INSERT INTO oauth2_users").
		WithArgs(o.ID, o.UserID, o.Provider, o.ProviderID, o.LinkedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), o)
	require.NoError(t, err)
}
func TestOAuth2UserRepository_CreateWithUser(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	user := &zmodel.User{
		ID:        "u123",
		Email:     "test@example.com",
		Name:      "Tester",
		CreatedAt: time.Now(),
	}
	oauthUser := &zmodel.OAuth2User{
		ID:         "oauth123",
		User:       user,
		Provider:   "google",
		ProviderID: "g_abc",
		LinkedAt:   time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Email, user.Name, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO oauth2_users").
		WithArgs(oauthUser.ID, user.ID, oauthUser.Provider, oauthUser.ProviderID, oauthUser.LinkedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	oauthUser.UserID = "" // será preenchido pela lógica
	created, err := repo.CreateWithUser(context.Background(), oauthUser)
	require.NoError(t, err)
	require.Equal(t, user.ID, created.UserID)
}

func TestOAuth2UserRepository_Update(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	o := &zmodel.OAuth2User{
		ID:         "abc",
		UserID:     "u42",
		Provider:   "github",
		ProviderID: "gh_abc",
		LinkedAt:   time.Now(),
	}

	mock.ExpectExec("UPDATE oauth2_users").
		WithArgs(o.ID, o.UserID, o.Provider, o.ProviderID, o.LinkedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), o)
	require.NoError(t, err)
}

func TestOAuth2UserRepository_Delete(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	mock.ExpectExec("DELETE FROM oauth2_users").
		WithArgs("abc").WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(context.Background(), "abc")
	require.NoError(t, err)
}

func TestOAuth2UserRepository_Transaction(t *testing.T) {
	repo, mock, teardown := setupOAuth2(t)
	defer teardown()

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := repo.Transaction(context.Background(), func(ctx context.Context) error {
		// lógica simulada dentro da tx
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
