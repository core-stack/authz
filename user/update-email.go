package user

import (
	"context"
	"time"

	"github.com/core-stack/authz/session"
	"github.com/core-stack/authz/zmodel"
	"github.com/google/uuid"
)

func (s *UserService) StartUpdateEmail(ctx context.Context, session *session.Session) error {
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}
	code := zmodel.Code{
		UserID:    user.ID,
		Token:     uuid.NewString(),
		Type:      zmodel.UpdateEmail,
		ExpiresAt: time.Now().Add(s.defaultCodeDuration),
	}
	err = s.codeRepo.Create(ctx, &code)
	if err != nil {
		return err
	}

	return s.emailSender.SendStartChangeEmail(ctx, *user, code.Token)
}

func (s *UserService) FinishUpdateEmail(ctx context.Context, session *session.Session, code string, email string) error {
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}
	user.Email = email
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	// TODO - remove oauth2 user relations
	return s.emailSender.SendFinishChangeEmail(ctx, *user)
}
