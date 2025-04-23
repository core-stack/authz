package auth

import (
	"time"

	"github.com/core-stack/authz/auth/oauth2"
	"github.com/core-stack/authz/email"
	"github.com/core-stack/authz/session"
	"github.com/core-stack/authz/zrepository"
)

type Auth struct {
	userRepo       zrepository.IUserRepository
	codeRepo       zrepository.ICodeRepository
	emailSender    email.ISender
	sessionManager session.ISessionManager

	defaultRoleId       int
	defaultCodeDuration time.Duration

	OAuth2 *oauth2.OAuth2
}

func NewAuth(
	userRepo zrepository.IUserRepository,
	codeRepo zrepository.ICodeRepository,
	emailSender email.ISender,
	oauth2 *oauth2.OAuth2,
) *Auth {
	return &Auth{
		userRepo:    userRepo,
		codeRepo:    codeRepo,
		emailSender: emailSender,
		OAuth2:      oauth2,
	}
}
