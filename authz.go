package authz

import (
	"log/slog"

	"github.com/core-stack/authz/auth"
	"github.com/core-stack/authz/auth/oauth2"
	"github.com/core-stack/authz/email"
	"github.com/core-stack/authz/jwt"
	"github.com/core-stack/authz/session"
	"github.com/core-stack/authz/store"
	"github.com/core-stack/authz/user"
	"github.com/core-stack/authz/zrepository"
)

type Authz struct {
	Auth              *auth.Auth
	User              *user.UserService
	Session           session.ISessionManager
	JWT               *jwt.JWTService
	EmailSender       email.ISender
	PermissionService *PermissionService
}

func New(opts ...OptionFunc) *Authz {
	initLogs()

	options := &AuthzOptions{}

	for _, opt := range opts {
		opt(options)
	}

	if options.db == nil {
		slog.Warn("Database connection not provided")
	}
	if options.emailSender == nil {
		slog.Warn("Email sender not provided")
	}

	if options.jwtSecret == "" {
		slog.Warn("JWT secret not provided, generating random secret...")
		options.jwtSecret = jwt.GenerateSecret(32)
	}

	if options.accessTokenDuration == 0 {
		options.accessTokenDuration = jwt.DefaultAccessTokenDuration
	}
	if options.refreshTokenDuration == 0 {
		options.refreshTokenDuration = jwt.DefaultRefreshTokenDuration
	}

	if options.sessionStore == nil {
		slog.Warn("Session store not provided, using in-memory store...")
		options.sessionStore = store.NewMemoryStore[session.Session]()
	}

	if options.sessionDuration == 0 {
		options.sessionDuration = session.DefaultSessionDuration
	}

	if options.providers == nil {
		options.providers = []oauth2.Provider{}
	}
	permissionService := NewPermissionService()
	jwtService := jwt.NewJWTService(
		options.jwtSecret,
		options.accessTokenDuration,
		options.refreshTokenDuration,
	)
	sessionManager := session.NewSessionManager(
		options.sessionStore,
		options.sessionDuration,
		jwtService,
	)
	emailSender := options.emailSender
	if emailSender == nil {
		emailSender = email.NewDummySender()
	}
	var (
		userRepo   zrepository.IUserRepository
		codeRepo   zrepository.ICodeRepository
		oauth2Repo zrepository.IOauth2UserRepository
	)

	if options.db != nil {
		userRepo = zrepository.NewUserRepository(options.db)
		codeRepo = zrepository.NewCodeRepository(options.db)
		oauth2Repo = zrepository.NewOAuth2UserRepository(options.db)
	}

	var userService *user.UserService
	if userRepo != nil && codeRepo != nil {
		userService = user.NewUserService(userRepo, codeRepo, emailSender, options.codeDuration)
	} else {
		slog.Warn("User service not created, user repository or code repository not provided")
	}

	var oauth2Service *oauth2.OAuth2
	if userRepo != nil && oauth2Repo != nil {
		oauth2Service = oauth2.NewOAuth2(userRepo, oauth2Repo, options.providers)
	} else {
		slog.Warn("OAuth2 service not created, user repository or oauth2 repository not provided")
	}

	var authService *auth.Auth
	if userRepo != nil && codeRepo != nil {
		authService = auth.NewAuth(userRepo, codeRepo, emailSender, oauth2Service)
	} else {
		slog.Warn("Auth service not created, user repository or code repository not provided")
	}
	authz := &Authz{
		PermissionService: permissionService,
		JWT:               jwtService,
		Session:           sessionManager,
		EmailSender:       emailSender,
		Auth:              authService,
		User:              userService,
	}
	return authz
}
