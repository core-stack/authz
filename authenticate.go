package authz

import (
	"context"
	"errors"
	"strings"

	"github.com/core-stack/authz/jwt"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
)

func (a *Authz) Authenticate(ctx context.Context, tokenExtractor func() (string, error)) (*jwt.TokenClaims, error) {
	token, err := tokenExtractor()
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, ErrNotAuthenticated
	}
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	claims, err := a.JWT.ParseToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
