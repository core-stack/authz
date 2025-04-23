package user

import (
	"context"

	"github.com/core-stack/authz/crypt"
	"github.com/core-stack/authz/session"
)

func (s *UserService) Delete(ctx context.Context, session *session.Session, password string) error {
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}
	if !crypt.ComparePassword(user.Password, password) {
		return ErrInvalidPassword
	}
	email := user.Email
	err = s.userRepo.Delete(ctx, user.ID)
	if err != nil {
		return err
	}
	return s.emailSender.SendDeleteAccount(ctx, email)
}
