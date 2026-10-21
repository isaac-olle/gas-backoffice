package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/key_factory"
)

type JwtTokenGenerator struct{}

func NewJwtTokenGenerator() *JwtTokenGenerator {
	return &JwtTokenGenerator{}
}

func (g *JwtTokenGenerator) Generate(_ context.Context, iss string, spec key_pairs.TokenSpec, kp *key_pairs.KeyPair, expiresIn time.Duration) (string, error) {
	p := kp.ToPrimitives()

	signingKey, err := key_factory.NewSigningKeyFromPEM([]byte(p.PrivateKey), p.Algorithm)
	if err != nil {
		return "", fmt.Errorf("invalid private key in key pair: %w", err)
	}

	header, err := jwt.NewTokenHeaderBuilder(signingKey, p.Id).Build(spec.IncludeJWK, spec.IncludeX5c)
	if err != nil {
		return "", fmt.Errorf("build token header: %w", err)
	}

	claims := jwt.NewTokenClaimBuilder(iss, spec.Sub, spec.Aud, spec.JTI, spec.ExtraClaims).Build(expiresIn)

	token, err := jwt.NewTokenSigner(signingKey).Sign(header, claims)
	if err != nil {
		return "", fmt.Errorf("token generation failed: %w", err)
	}

	return token, nil
}
