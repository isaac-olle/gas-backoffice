package infrastructure

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomCSRFTokenGenerator struct{}

func (RandomCSRFTokenGenerator) GenerateCSRFToken() string {
	var token [32]byte
	_, _ = rand.Read(token[:])
	return base64.RawURLEncoding.EncodeToString(token[:])
}
