package key_pairs_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
)

// newKeyPairMother creates a valid KeyPair with random IDs for use in tests.
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
