package auth

import (
	"context"
	"errors"

	"github.com/core-stack/authz/crypt"
	"github.com/core-stack/authz/jwt"
	"github.com/core-stack/authz/zmodel"
)

var (
	ErrInvalidLogin       = errors.New("invalid login")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginOpts struct {
	Email    string
	Password string
	Name     string

	UserAgent string
	IP        string
}

type LoginOption func(*LoginOpts)

func WithEmail(email, password string) LoginOption {
	return func(opts *LoginOpts) {
		opts.Email = email
		opts.Password = password
	}
}
func WithUsername(name, password string) LoginOption {
	return func(opts *LoginOpts) {
		opts.Name = name
		opts.Password = password
	}
}
func WithUserAgent(ua string) LoginOption {
	return func(opts *LoginOpts) {
		opts.UserAgent = ua
	}
}
func WithIP(ip string) LoginOption {
	return func(opts *LoginOpts) {
		opts.IP = ip
	}
}

func (s *Auth) Login(ctx context.Context, opts ...LoginOption) (*jwt.TokenPair, error) {
	loginOpts := LoginOpts{}
	for _, opt := range opts {
		opt(&loginOpts)
	}
	var err error
	var user *zmodel.User
	if loginOpts.Email != "" {
		user, err = s.userRepo.GetByEmail(ctx, loginOpts.Email)
		if err != nil {
			return nil, err
		}
	} else if loginOpts.Name != "" {
		user, err = s.userRepo.GetByUsername(ctx, loginOpts.Name)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrInvalidLogin
	}

	if crypt.ComparePassword(user.Password, loginOpts.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.sessionManager.CreateSession(ctx, *user, map[string]string{
		"user_agent": loginOpts.UserAgent,
		"ip":         loginOpts.IP,
	})
}
