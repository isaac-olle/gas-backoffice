package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/infrastructure"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresSigningKeyRepositorySearchActive(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, protectedResourceID := integrationtest.ProtectedResource(t, pool)
	keyID := uuid.NewString()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM security.key_pairs WHERE id=$1`, keyID) })
	_, err := pool.Exec(ctx, `INSERT INTO security.key_pairs (id,protected_resource_id,algorithm,public_key,private_key) VALUES ($1,$2,$3,$4,$5)`, keyID, protectedResourceID, "RS256", []byte("public-key"), []byte("private-key"))
	require.NoError(t, err)

	keys, err := infrastructure.NewPostgresSigningKeyRepository(pool).SearchActive(ctx)
	require.NoError(t, err)
	var found bool
	for _, key := range keys {
		if key.Kid == keyID {
			found = true
			assert.Equal(t, "RS256", key.Algorithm)
			assert.Equal(t, []byte("public-key"), key.PublicKeyPEM)
		}
	}
	assert.True(t, found)

	_, err = pool.Exec(ctx, `UPDATE backoffice.protected_resources SET enabled=FALSE WHERE id=$1`, protectedResourceID)
	require.NoError(t, err)
	keys, err = infrastructure.NewPostgresSigningKeyRepository(pool).SearchActive(ctx)
	require.NoError(t, err)
	for _, key := range keys {
		assert.NotEqual(t, keyID, key.Kid)
	}
}
