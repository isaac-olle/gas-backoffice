package generate_jwt_by_pr_id_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
)

type mockTokenGenerator struct {
	token      string
	err        error
	gotKeyPair *key_pairs.KeyPair
	gotIss     string
	gotSpec    key_pairs.TokenSpec
}

func (m *mockTokenGenerator) Generate(_ context.Context, iss string, spec key_pairs.TokenSpec, kp *key_pairs.KeyPair, _ time.Duration) (string, error) {
	m.gotKeyPair = kp
	m.gotIss = iss
	m.gotSpec = spec
	return m.token, m.err
}

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
