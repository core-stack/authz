package authz

import (
	"database/sql"
	"time"

	"github.com/core-stack/authz/auth/oauth2"
	"github.com/core-stack/authz/email"
	"github.com/core-stack/authz/session"
	"github.com/core-stack/authz/store"
)

type AuthzOptions struct {
	db          *sql.DB
	emailSender email.ISender

	providers []oauth2.Provider

	sessionStore    store.Store[session.Session]
	sessionDuration time.Duration

	jwtSecret                                 string
	accessTokenDuration, refreshTokenDuration time.Duration

	codeDuration time.Duration
}

type OptionFunc func(*AuthzOptions)

func WithDB(db *sql.DB) OptionFunc {
	return func(o *AuthzOptions) {
		o.db = db
	}
}

func WithEmailSender(sender email.ISender) OptionFunc {
	return func(o *AuthzOptions) {
		o.emailSender = sender
	}
}

func WithOAuth2Providers(providers ...oauth2.Provider) OptionFunc {
	return func(o *AuthzOptions) {
		o.providers = providers
	}
}

func WithSessionStore(store store.Store[session.Session]) OptionFunc {
	return func(o *AuthzOptions) {
		o.sessionStore = store
	}
}

func WithSessionDuration(d time.Duration) OptionFunc {
	return func(o *AuthzOptions) {
		o.sessionDuration = d
	}
}

func WithCodeDuration(d time.Duration) OptionFunc {
	return func(o *AuthzOptions) {
		o.codeDuration = d
	}
}

func WithJWTSettings(secret string, accessTokenDuration, refreshTokenDuration time.Duration) func(*AuthzOptions) {
	return func(o *AuthzOptions) {
		o.jwtSecret = secret
		o.accessTokenDuration = accessTokenDuration
		o.refreshTokenDuration = refreshTokenDuration
	}
}
