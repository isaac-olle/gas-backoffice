package jwks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks"
)

func TestSigningKeyProjection_ToJWK(t *testing.T) {
	t.Parallel()

	t.Run("renders the public key with kid and use", func(t *testing.T) {
		t.Parallel()
		jwk, err := jwks.NewSigningKeyProjection("kid-1", "RS256", aPublicKeyPEM(t)).ToJWK()
		require.NoError(t, err)
		assert.Equal(t, "kid-1", jwk["kid"])
		assert.Equal(t, "sig", jwk["use"])
		assert.Equal(t, "RSA", jwk["kty"])
		assert.Equal(t, "RS256", jwk["alg"])
		assert.NotEmpty(t, jwk["n"])
		assert.NotEmpty(t, jwk["e"])
		assert.NotContains(t, jwk, "d")
	})

	t.Run("returns error for an invalid pem", func(t *testing.T) {
		t.Parallel()
		_, err := jwks.NewSigningKeyProjection("kid-1", "RS256", []byte("not a pem")).ToJWK()
		assert.Error(t, err)
	})

	t.Run("returns error for an unsupported algorithm", func(t *testing.T) {
		t.Parallel()
		_, err := jwks.NewSigningKeyProjection("kid-1", "HS256", aPublicKeyPEM(t)).ToJWK()
		assert.Error(t, err)
	})
}
