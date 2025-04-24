package oauth2

import "errors"

var (
	ErrProviderNotFound           = errors.New("provider not found")
	ErrProviderNotSupportUserInfo = errors.New("provider does not support user info")
)
