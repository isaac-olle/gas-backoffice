package infrastructure_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

type integrationEncrypter struct{}

func (integrationEncrypter) Encrypt(plain string) (string, error) { return "encrypted:" + plain, nil }
func (integrationEncrypter) Decrypt(ciphertext string) (string, error) {
	const prefix = "encrypted:"
	if len(ciphertext) < len(prefix) || ciphertext[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid encrypted value")
	}
	return ciphertext[len(prefix):], nil
}
func (integrationEncrypter) Seal(value []byte) ([]byte, error) { return value, nil }
func (integrationEncrypter) Open(value []byte) ([]byte, error) { return value, nil }

func TestPostgresKeyPairRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawProtectedResourceID := integrationtest.ProtectedResource(t, pool)
	protectedResourceID, err := ids.NewProtectedResourceIdFromString(rawProtectedResourceID)
	require.NoError(t, err)
	keyID := ids.NewKeyPairId()
	keyGenerator := infrastructure.NewPEMKeyGenerator()
	privateKey, publicKey, err := keyGenerator.Generate("ES256")
	require.NoError(t, err)
	keyPair, err := key_pairs.NewKeyPair(keyID.String(), protectedResourceID.String(), "ES256", privateKey, publicKey)
	require.NoError(t, err)
	repo := infrastructure.NewPostgresKeyPairRepository(pool, integrationEncrypter{})
	t.Cleanup(func() { _ = repo.HardDelete(ctx, keyID) })
	require.NoError(t, repo.Save(ctx, keyPair))

	byID, err := repo.SearchById(ctx, keyID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.Equal(t, []byte(privateKey), byID.ToPrimitives().PrivateKey)
	assert.Equal(t, []byte(publicKey), byID.ToPrimitives().PublicKey)

	byResource, err := repo.SearchByProtectedResourceId(ctx, protectedResourceID)
	require.NoError(t, err)
	require.NotNil(t, byResource)
	assert.Equal(t, keyID.String(), byResource.ToPrimitives().Id)

	var storedPrivateKey []byte
	require.NoError(t, pool.QueryRow(ctx, `SELECT private_key FROM security.key_pairs WHERE id=$1`, keyID.String()).Scan(&storedPrivateKey))
	assert.NotEqual(t, []byte(privateKey), storedPrivateKey)

	require.NoError(t, repo.SoftDelete(ctx, keyID))
	deleted, err := repo.SearchById(ctx, keyID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
