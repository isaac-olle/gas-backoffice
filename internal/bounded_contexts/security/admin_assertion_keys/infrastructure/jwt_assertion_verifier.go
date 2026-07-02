package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/jwt"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/admin_assertion_keys"
)

type JwtAssertionVerifier struct {
	verifier jwt.TokenVerifier
}

func NewJwtAssertionVerifier(verifier jwt.TokenVerifier) *JwtAssertionVerifier {
	return &JwtAssertionVerifier{verifier: verifier}
}

func (v *JwtAssertionVerifier) Verify(_ context.Context, assertion string, key *admin_assertion_keys.AdminAssertionKey) (map[string]any, error) {
	p := key.ToPrimitives()

	voKey, err := jwt.NewKey([]byte(p.PublicKey), p.Algorithm, p.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid public key in assertion key: %w", err)
	}

	claims, err := v.verifier.Verify(assertion, voKey)
	if err != nil {
		return nil, fmt.Errorf("assertion signature verification failed: %w", err)
	}

	return claims, nil
}
