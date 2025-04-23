package session

import (
	"context"
	"time"

	"github.com/core-stack/authz/jwt"
	zrole "github.com/core-stack/authz/role"
	"github.com/core-stack/authz/zmodel"
	"github.com/google/uuid"
)

func (s *SessionManager) CreateSession(ctx context.Context, user zmodel.User, extra map[string]string) (*jwt.TokenPair, error) {
	session := Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		RoleID:    user.RoleID,
		Extra:     extra,
		ExpiresAt: time.Now().Add(s.sessionDuration),
	}

	err := s.store.Set(ctx, session)
	if err != nil {
		return nil, err
	}
	role, err := zrole.GetRole(user.RoleID)
	if err != nil {
		return nil, err
	}
	tokenPair, err := s.jwtService.Generate(ctx, user.ID, session.ID, role.Permissions)
	return tokenPair, err
}
