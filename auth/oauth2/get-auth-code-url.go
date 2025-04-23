package oauth2

import "context"

func (s *OAuth2) GetAuthCodeURL(ctx context.Context, provider string, state string) (string, error) {
	p, err := s.getProvider(provider)
	if err != nil {
		return "", err
	}
	return p.AuthCodeURL(state), nil
}
