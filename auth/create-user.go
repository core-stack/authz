package auth

import (
	"context"
	"time"

	"github.com/core-stack/authz/zmodel"
	"github.com/google/uuid"
)

type CreateUserOpts struct {
	Nickname string
	Name     string
	Email    string
	Password string
	IsAdmin  bool
	RoleId   int
}
type CreateUserOption func(*CreateUserOpts)

func WithNickname(
	nick string,
) CreateUserOption {
	return func(opts *CreateUserOpts) {
		opts.Nickname = nick
	}
}
func WithRole(
	role int,
) CreateUserOption {
	return func(opts *CreateUserOpts) {
		opts.RoleId = role
	}
}
func WithIsAdmin(
	isAdmin bool,
) CreateUserOption {
	return func(opts *CreateUserOpts) {
		opts.IsAdmin = isAdmin
	}
}

func (s *Auth) CreateUser(ctx context.Context, name, email, password string, opts ...CreateUserOption) (*zmodel.User, error) {
	createUser := CreateUserOpts{Name: name, Email: email, Password: password, RoleId: s.defaultRoleId}

	for _, opt := range opts {
		opt(&createUser)
	}
	user := zmodel.User{
		Username: createUser.Nickname,
		Name:     createUser.Name,
		Email:    createUser.Email,
		Password: createUser.Password,
		RoleID:   createUser.RoleId,
		IsAdmin:  createUser.IsAdmin,
		Status:   zmodel.Pending,
	}
	code := zmodel.Code{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(s.defaultCodeDuration),
	}
	s.userRepo.Transaction(ctx, func(ctx context.Context) error {
		err := s.userRepo.Create(ctx, &user)
		if err != nil {
			return err
		}
		code.UserID = user.ID
		return s.codeRepo.Create(ctx, &code)
	})
	err := s.emailSender.SendActiveAccount(ctx, user, code.Token)
	return &user, err
}
