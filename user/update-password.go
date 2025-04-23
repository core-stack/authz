package user

import (
	"context"
	"errors"

	"github.com/core-stack/authz/crypt"
	"github.com/core-stack/authz/session"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

func (s *UserService) UpdatePassword(ctx context.Context, session *session.Session, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}
	if user.Password == "" && !crypt.ComparePassword(user.Password, oldPassword) {
		return ErrInvalidPassword
	}

	user.Password, err = crypt.HashPassword(newPassword)
	if err != nil {
		return err
	}
	err = s.userRepo.UpdatePassword(ctx, user.ID, user.Password)
	if err != nil {
		return err
	}
	return s.emailSender.SendChangePassword(ctx, *user)
}
