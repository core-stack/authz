package user

import (
	"context"

	"github.com/core-stack/authz/session"
	"github.com/core-stack/authz/zmodel"
)

func (s *UserService) Get(ctx context.Context, session *session.Session) (*zmodel.User, error) {
	return s.userRepo.GetByID(ctx, session.UserID)
}
