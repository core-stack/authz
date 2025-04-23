package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	DefaultAccessTokenDuration  = 15 * time.Minute
	DefaultRefreshTokenDuration = 7 * 24 * time.Hour
)

type JWTService struct {
	secretKey       []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

type TokenPair struct {
	AccessToken   string
	RefreshToken  string
	AccessExpiry  time.Time
	RefreshExpiry time.Time
}

type TokenClaims struct {
	UserID      string
	SessionID   string
	Permissions int
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessDur, refreshDur time.Duration) *JWTService {
	return &JWTService{
		secretKey:       []byte(secret),
		accessDuration:  accessDur,
		refreshDuration: refreshDur,
	}
}

func (j *JWTService) Generate(ctx context.Context, userID, sessionID string, permissions int) (*TokenPair, error) {
	now := time.Now()

	accessExp := now.Add(j.accessDuration)
	refreshExp := now.Add(j.refreshDuration)

	accessClaims := TokenClaims{
		UserID:      userID,
		SessionID:   sessionID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        sessionID,
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ID:        sessionID,
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExp,
		RefreshExpiry: refreshExp,
	}, nil
}

func (j *JWTService) ParseToken(_ context.Context, tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func (j *JWTService) Refresh(ctx context.Context, refreshToken string, permissions int) (*TokenPair, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	userID := claims.Subject
	sessionID := claims.ID

	return j.Generate(ctx, userID, sessionID, permissions)
}
