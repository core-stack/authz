package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/core-stack/authz/crypt"
)

func (s *Auth) ResetPassword(ctx context.Context, token string, password string) error {
	code, err := s.codeRepo.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if code.ExpiresAt.Before(time.Now()) {
		return ErrCodeExpired
	}
	if code.UsedAt.Valid {
		return ErrCodeUsed
	}
	user, err := s.userRepo.GetByID(ctx, code.UserID)
	if err != nil {
		return err
	}
	user.Password, err = crypt.HashPassword(password)
	if err != nil {
		return err
	}
	return s.userRepo.Transaction(ctx, func(ctx context.Context) error {
		err = s.userRepo.UpdatePassword(ctx, user.ID, user.Password)
		if err != nil {
			return err
		}
		code.UsedAt = sql.NullTime{Time: time.Now(), Valid: true}
		err = s.codeRepo.Update(ctx, code)
		if err != nil {
			return err
		}
		return s.emailSender.SendNotifyResetPassword(ctx, *user)
	})
}
