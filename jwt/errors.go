package jwt

import "errors"

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidClaims       = errors.New("invalid claims")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
