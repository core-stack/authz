package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/core-stack/authz/zmodel"
)

var (
	ErrCodeExpired = errors.New("code expired")
	ErrCodeUsed    = errors.New("code already used")
)

func (s *Auth) ActiveAccount(ctx context.Context, token string) error {
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
	user.Status = zmodel.Active
	return s.userRepo.Transaction(ctx, func(ctx context.Context) error {
		err = s.userRepo.Update(ctx, user)
		if err != nil {
			return err
		}
		code.UsedAt = sql.NullTime{Time: time.Now(), Valid: true}
		err = s.codeRepo.Update(ctx, code)
		if err != nil {
			return err
		}
		return nil
	})
}
