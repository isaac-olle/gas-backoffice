package infrastructure_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestJwtTokenAdapters(t *testing.T) {
	t.Parallel()
	privateKey, publicKey, err := infrastructure.NewPEMKeyGenerator().Generate("ES256")
	require.NoError(t, err)
	keyID := ids.NewKeyPairId()
	keyPair, err := key_pairs.NewKeyPair(keyID.String(), ids.NewProtectedResourceId().String(), "ES256", privateKey, publicKey)
	require.NoError(t, err)

	ctx := context.Background()
	token, err := infrastructure.NewJwtTokenGenerator().Generate(ctx, "issuer", key_pairs.TokenSpec{
		Sub: "user-123", Aud: "api", JTI: "token-123", ExtraClaims: map[string]any{"scope": "read"},
	}, keyPair, time.Minute)
	require.NoError(t, err)

	decodedKeyID, err := infrastructure.NewJwtDecoder().ExtractKeyId(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, keyID.String(), decodedKeyID)

	claims, err := infrastructure.NewJwtTokenVerifier().Verify(ctx, token, keyPair, map[string]any{
		"iss": "issuer", "sub": "user-123", "aud": "api", "scope": "read",
	}, true)
	require.NoError(t, err)
	assert.Equal(t, "read", claims["scope"])
	assert.Equal(t, "token-123", claims["jti"])

	_, err = infrastructure.NewJwtTokenVerifier().Verify(ctx, token, keyPair, map[string]any{"sub": "another-user"}, true)
	assert.Error(t, err)

	expiredToken, err := infrastructure.NewJwtTokenGenerator().Generate(ctx, "issuer", key_pairs.TokenSpec{Sub: "user-123", Aud: "api"}, keyPair, -time.Minute)
	require.NoError(t, err)
	_, err = infrastructure.NewJwtTokenVerifier().Verify(ctx, expiredToken, keyPair, nil, true)
	assert.Error(t, err)
}
