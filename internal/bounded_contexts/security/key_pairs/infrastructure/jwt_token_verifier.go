package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/jwt"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs"
)

type JwtTokenVerifier struct {
	verifier jwt.TokenVerifier
}

func NewJwtTokenVerifier(verifier jwt.TokenVerifier) *JwtTokenVerifier {
	return &JwtTokenVerifier{verifier: verifier}
}

func (v *JwtTokenVerifier) Verify(_ context.Context, token string, kp *key_pairs.KeyPair) (*key_pairs.TokenClaims, error) {
	p := kp.ToPrimitives()

	jwtKey, err := jwt.NewKey([]byte(p.PublicKey), p.Algorithm, p.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid public key in key pair: %w", err)
	}

	rawClaims, err := v.verifier.Verify(token, jwtKey)
	if err != nil {
		return nil, fmt.Errorf("token signature verification failed: %w", err)
	}

	claims, err := key_pairs.NewTokenClaims(rawClaims)
	if err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	return claims, nil
}
