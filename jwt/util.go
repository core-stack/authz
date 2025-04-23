package jwt

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecret(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
