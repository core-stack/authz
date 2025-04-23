package oauth2

import (
	"context"
	"database/sql"
	"errors"

	"github.com/core-stack/authz/jwt"
	"github.com/core-stack/authz/zmodel"
)

type CallbackOptions struct {
	UserAgent string
	IP        string
}
type CallbackOption func(*CallbackOptions)

func WithUserAgent(ua string) CallbackOption {
	return func(opts *CallbackOptions) {
		opts.UserAgent = ua
	}
}

func WithIP(ip string) CallbackOption {
	return func(opts *CallbackOptions) {
		opts.IP = ip
	}
}
func (s *OAuth2) Callback(ctx context.Context, provider string, code string, opts ...CallbackOption) (*jwt.TokenPair, error) {
	options := CallbackOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	p, err := s.getProvider(provider)
	if err != nil {
		return nil, err
	}
	token, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}
	user, err := p.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}
	oauth, err := s.oauth2Repo.FindByProviderID(ctx, provider, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if oauth == nil {
		oauth, err = s.oauth2Repo.CreateWithUser(ctx, &zmodel.OAuth2User{
			Provider:   provider,
			ProviderID: user.ID,
			User: &zmodel.User{
				Name:   user.Name,
				Email:  user.Email,
				Status: zmodel.Active,
				RoleID: s.defaultRoleId,
			},
		})
	}
	if err != nil {
		return nil, err
	}
	return s.sessionManager.CreateSession(
		ctx,
		*oauth.User,
		map[string]string{
			"user_agent": options.UserAgent,
			"ip":         options.IP,
		},
	)
}
