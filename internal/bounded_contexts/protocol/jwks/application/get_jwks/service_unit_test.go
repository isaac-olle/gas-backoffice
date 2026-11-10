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

	tests := []struct {
		name     string
		repo     func(t *testing.T) *stubRepository
		wantErr  bool
		wantKids []string
	}{
		{
			name: "returns one jwk per active signing key",
			repo: func(t *testing.T) *stubRepository {
				return &stubRepository{keys: []*jwks.SigningKeyProjection{
					jwks.NewSigningKeyProjection("kid-1", "RS256", aPublicKeyPEM(t)),
					jwks.NewSigningKeyProjection("kid-2", "RS256", aPublicKeyPEM(t)),
				}}
			},
			wantKids: []string{"kid-1", "kid-2"},
		},
		{
			name: "returns an empty set when there are no keys",
			repo: func(t *testing.T) *stubRepository { return &stubRepository{} },
		},
		{
			name:    "returns error when repo errors",
			repo:    func(t *testing.T) *stubRepository { return &stubRepository{err: fmt.Errorf("db error")} },
			wantErr: true,
		},
		{
			name: "returns error when a key cannot be rendered",
			repo: func(t *testing.T) *stubRepository {
				return &stubRepository{keys: []*jwks.SigningKeyProjection{
					jwks.NewSigningKeyProjection("kid-1", "RS256", []byte("broken")),
				}}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			set, err := newService(tt.repo(t)).Get(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantKids == nil {
				assert.Empty(t, set)
				return
			}
			require.Len(t, set, len(tt.wantKids))
			for i, kid := range tt.wantKids {
				assert.Equal(t, kid, set[i]["kid"])
			}
		})
	}
}
