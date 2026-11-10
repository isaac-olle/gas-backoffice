package verify_jwt_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
)

// mockTokenVerifier stubs key_pairs.TokenVerifier.
type mockTokenVerifier struct {
	claims              map[string]any
	err                 error
	gotToken            string
	gotKeyPair          *key_pairs.KeyPair
	gotExpectedClaims   map[string]any
	gotVerifyExpiration bool
}

func (m *mockTokenVerifier) Verify(_ context.Context, token string, kp *key_pairs.KeyPair, expectedClaims map[string]any, verifyExpiration bool) (map[string]any, error) {
	m.gotToken = token
	m.gotKeyPair = kp
	m.gotExpectedClaims = expectedClaims
	m.gotVerifyExpiration = verifyExpiration
	return m.claims, m.err
}

// mockJwtDecoder stubs key_pairs.JwtDecoder.
type mockJwtDecoder struct {
	kid      string
	err      error
	gotToken string
}

func (m *mockJwtDecoder) ExtractKeyId(_ context.Context, token string) (string, error) {
	m.gotToken = token
	return m.kid, m.err
}

// newKeyPairMother builds a valid KeyPair with fresh random IDs.
func newKeyPairMother(t *testing.T) *key_pairs.KeyPair {
	t.Helper()
	privateKeyPEM, publicKeyPEM, err := keyPairInfra.NewPEMKeyGenerator().Generate("RS256")
	require.NoError(t, err)
	kp, err := key_pairs.NewKeyPair(
		uuid.New().String(),
		uuid.New().String(),
		"RS256",
		privateKeyPEM,
		publicKeyPEM,
	)
	require.NoError(t, err)
	return kp
}

// newValidClaimsMother builds a valid claims map.
func newValidClaimsMother(issuer string) map[string]any {
	return map[string]any{
		"sub":  uuid.New().String(),
		"iss":  issuer,
		"aud":  issuer,
		"role": "admin",
	}
}
