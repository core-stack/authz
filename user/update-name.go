package user

import (
	"context"

	"github.com/core-stack/authz/session"
)

func (s *UserService) UpdateName(ctx context.Context, session *session.Session, name string) error {
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}
	user.Name = name
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) UpdateUsername(ctx context.Context, session *session.Session, username string) error {
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}
	user.Username = username
	return s.userRepo.Update(ctx, user)
}
