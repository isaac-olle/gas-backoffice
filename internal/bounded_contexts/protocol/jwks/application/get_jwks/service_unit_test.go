package get_jwks_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/application/get_jwks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/domain_services"
)

type stubRepository struct {
	keys []*jwks.SigningKeyProjection
	err  error
}

func (s *stubRepository) SearchActive(context.Context) ([]*jwks.SigningKeyProjection, error) {
	return s.keys, s.err
}

func aPublicKeyPEM(t *testing.T) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func newService(repo jwks.SigningKeyRepository) *get_jwks.Service {
	return get_jwks.NewService(domain_services.NewSigningKeyFinder(repo))
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns one jwk per active signing key", func(t *testing.T) {
		t.Parallel()
		repo := &stubRepository{keys: []*jwks.SigningKeyProjection{
			jwks.NewSigningKeyProjection("kid-1", "RS256", aPublicKeyPEM(t)),
			jwks.NewSigningKeyProjection("kid-2", "RS256", aPublicKeyPEM(t)),
		}}

		set, err := newService(repo).Get(context.Background())
		require.NoError(t, err)
		require.Len(t, set, 2)
		assert.Equal(t, "kid-1", set[0]["kid"])
		assert.Equal(t, "kid-2", set[1]["kid"])
	})

	t.Run("returns an empty set when there are no keys", func(t *testing.T) {
		t.Parallel()
		set, err := newService(&stubRepository{}).Get(context.Background())
		require.NoError(t, err)
		assert.Empty(t, set)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		_, err := newService(&stubRepository{err: fmt.Errorf("db error")}).Get(context.Background())
		assert.Error(t, err)
	})

	t.Run("returns error when a key cannot be rendered", func(t *testing.T) {
		t.Parallel()
		repo := &stubRepository{keys: []*jwks.SigningKeyProjection{
			jwks.NewSigningKeyProjection("kid-1", "RS256", []byte("broken")),
		}}
		_, err := newService(repo).Get(context.Background())
		assert.Error(t, err)
	})
}
