package auth

import (
	"context"
	"time"

	"github.com/core-stack/authz/zmodel"
	"github.com/google/uuid"
)

func (s *Auth) ForgetPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	code := zmodel.Code{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(s.defaultCodeDuration),
		UserID:    user.ID,
	}
	err = s.codeRepo.Create(ctx, &code)
	if err != nil {
		return err
	}
	return s.emailSender.SendResetPassword(ctx, *user, code.Token)
}
