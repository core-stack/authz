package oauth2

import (
	"errors"

	"github.com/core-stack/authz/session"
	"github.com/core-stack/authz/zrepository"
)

type OAuth2 struct {
	providers []Provider

	// repositories
	userRepo   zrepository.IUserRepository
	oauth2Repo zrepository.IOauth2UserRepository

	defaultRoleId int

	sessionManager session.ISessionManager
}

func NewOAuth2(
	userRepo zrepository.IUserRepository,
	oauth2Repo zrepository.IOauth2UserRepository,
	providers []Provider,
) *OAuth2 {
	return &OAuth2{
		userRepo:   userRepo,
		oauth2Repo: oauth2Repo,
		providers:  providers,
	}
}

func (s *OAuth2) getProvider(provider string) (Provider, error) {
	for _, p := range s.providers {
		if p.Name == provider {
			return p, nil
		}
	}
	return Provider{}, errors.New("provider not found")
}
