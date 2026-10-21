package infrastructure

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure/jwt_verifiers"
)

type JwtTokenVerifier struct {
	signatureVerifier     *jwt_verifiers.Signature
	expectedClaimVerifier *jwt_verifiers.ExpectedClaim
	expirationVerifier    *jwt_verifiers.ExpirationVerifier
}

func NewJwtTokenVerifier() *JwtTokenVerifier {
	return &JwtTokenVerifier{
		signatureVerifier:     jwt_verifiers.NewSignatureVerifier(),
		expectedClaimVerifier: jwt_verifiers.NewExpectedClaimVerifier(),
		expirationVerifier:    jwt_verifiers.NewExpirationVerifier(),
	}
}

func (v *JwtTokenVerifier) Verify(
	_ context.Context,
	token string,
	kp *key_pairs.KeyPair,
	expectedClaims map[string]any,
	verifyExpiration bool,
) (map[string]any, error) {
	claims, err := v.signatureVerifier.Verify(token, kp)
	if err != nil {
		return nil, err
	}

	if verifyExpiration {
		if err := v.expirationVerifier.Verify(claims); err != nil {
			return nil, err
		}
	}

	if err := v.expectedClaimVerifier.Verify(claims, expectedClaims); err != nil {
		return nil, err
	}

	return claims, nil
}
