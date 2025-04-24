package email

import (
	"context"

	"github.com/core-stack/authz/zmodel"
)

type ISender interface {
	Send(ctx context.Context, to, subject, body string) error
	SendActiveAccount(ctx context.Context, user zmodel.User, code string) error
	SendResetPassword(ctx context.Context, user zmodel.User, code string) error
	SendNotifyResetPassword(ctx context.Context, user zmodel.User) error
	SendChangePassword(ctx context.Context, user zmodel.User) error
	SendDeleteAccount(ctx context.Context, email string) error
}
