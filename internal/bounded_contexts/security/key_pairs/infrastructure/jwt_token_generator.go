package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/jwt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
)

type JwtTokenGenerator struct {
	generator jwt.TokenGenerator
}

func NewJwtTokenGenerator(generator jwt.TokenGenerator) *JwtTokenGenerator {
	return &JwtTokenGenerator{generator: generator}
}

func (g *JwtTokenGenerator) Generate(_ context.Context, claims map[string]any, kp *key_pairs.KeyPair, expiresIn time.Duration) (string, error) {
	p := kp.ToPrimitives()

	jwtKey, err := jwt.NewKey([]byte(p.PrivateKey), p.Algorithm, p.Id)
	if err != nil {
		return "", fmt.Errorf("invalid private key in key pair: %w", err)
	}

	token, err := g.generator.Generate(claims, jwtKey, expiresIn)
	if err != nil {
		return "", fmt.Errorf("token generation failed: %w", err)
	}

	return token, nil
}
