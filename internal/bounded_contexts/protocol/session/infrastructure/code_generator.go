package infrastructure

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomCodeGenerator struct{}

func NewRandomCodeGenerator() *RandomCodeGenerator {
	return &RandomCodeGenerator{}
}

func (g *RandomCodeGenerator) Generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
