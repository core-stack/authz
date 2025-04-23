package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type UserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

type Provider struct {
	Name            string
	AuthURL         string
	TokenURL        string
	UserInfoURL     string
	Scopes          []string
	ClientID        string
	ClientSecret    string
	RedirectURI     string
	ExtractUserInfo func([]byte) (UserInfo, error)
}

func (p *Provider) AuthCodeURL(state string) string {
	values := url.Values{}
	values.Set("client_id", p.ClientID)
	values.Set("redirect_uri", p.RedirectURI)
	values.Set("response_type", "code")
	values.Set("scope", strings.Join(p.Scopes, " "))
	values.Set("state", state)

	return fmt.Sprintf("%s?%s", p.AuthURL, values.Encode())
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Expiry       int    `json:"expires_in"`
	IDToken      string `json:"id_token,omitempty"`
}

func (p *Provider) ExchangeCode(ctx context.Context, code string) (TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", p.ClientID)
	data.Set("client_secret", p.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", p.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return TokenResponse{}, fmt.Errorf("token exchange failed: %s", body)
	}

	var tr TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tr)
	return tr, err
}

func (p *Provider) GetUserInfo(ctx context.Context, token string) (UserInfo, error) {
	if p.UserInfoURL == "" {
		return UserInfo{}, errors.New("provider does not support user info")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.UserInfoURL, nil)
	if err != nil {
		return UserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return UserInfo{}, fmt.Errorf("user info failed: %s", body)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, err
	}
	return p.ExtractUserInfo(b)
}
