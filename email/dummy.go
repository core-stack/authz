package email

import (
	"context"
	"log/slog"
	"os"

	"github.com/core-stack/authz/zmodel"
)

type DummySender struct {
	logger *slog.Logger
}

func NewDummySender() ISender {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &DummySender{
		logger: logger,
	}
}
func (s *DummySender) Send(ctx context.Context, to, subject, body string) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("Sending email", "to", to, "subject", subject, "body", body)
	return nil
}

func (s *DummySender) SendActiveAccount(ctx context.Context, user zmodel.User, code string) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendActiveAccount",
		"user_id", user.ID,
		"email", user.Email,
		"code", code,
	)
	return nil
}

func (s *DummySender) SendResetPassword(ctx context.Context, user zmodel.User, code string) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendResetPassword",
		"user_id", user.ID,
		"email", user.Email,
		"code", code,
	)
	return nil
}

func (s *DummySender) SendNotifyResetPassword(ctx context.Context, user zmodel.User) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendNotifyResetPassword",
		"user_id", user.ID,
		"email", user.Email,
	)
	return nil
}

func (s *DummySender) SendChangePassword(ctx context.Context, user zmodel.User) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendChangePassword",
		"user_id", user.ID,
		"email", user.Email,
	)
	return nil
}

func (s *DummySender) SendStartChangeEmail(ctx context.Context, user zmodel.User, code string) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendStartChangeEmail",
		"user_id", user.ID,
		"old_email", user.Email,
		"code", code,
	)
	return nil
}

func (s *DummySender) SendFinishChangeEmail(ctx context.Context, user zmodel.User) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendFinishChangeEmail",
		"user_id", user.ID,
		"new_email", user.Email,
	)
	return nil
}

func (s *DummySender) SendDeleteAccount(ctx context.Context, user zmodel.User) error {
	slog.Info("Sending email with dummy sender")
	s.logger.Info("SendDeleteAccount",
		"email", user.Email,
	)
	return nil
}
