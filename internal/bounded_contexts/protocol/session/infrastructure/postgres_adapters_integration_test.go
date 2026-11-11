package infrastructure_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresClientRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, subscriptionID := integrationtest.BackofficeSubscription(t, pool)
	protectedResourceID, scopeID, clientID := uuid.NewString(), uuid.NewString(), uuid.NewString()
	_, err := pool.Exec(ctx, `INSERT INTO backoffice.protected_resources (id,name,subscription_id) VALUES ($1,$2,$3)`, protectedResourceID, "Protocol Resource", subscriptionID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO backoffice.clients (id,name,subscription_id,allowed_grant_types,redirect_uris,public) VALUES ($1,$2,$3,$4,$5,$6)`, clientID, "Protocol Client", subscriptionID, []string{"authorization_code", "client_credentials"}, []string{"https://app.example/callback"}, true)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO backoffice.scopes (id,name,protected_resource_id,subscription_id,allowed_client_ids) VALUES ($1,$2,$3,$4,$5)`, scopeID, "profile", protectedResourceID, subscriptionID, []string{clientID})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.scopes WHERE id=$1`, scopeID)
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.clients WHERE id=$1`, clientID)
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.protected_resources WHERE id=$1`, protectedResourceID)
	})

	repo := infrastructure.NewPostgresClientRepository(pool)
	client, err := repo.FindById(ctx, clientID)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "Protocol Client", client.Name)
	assert.Equal(t, []string{"authorization_code", "client_credentials"}, client.AllowedGrantTypes)
	require.Len(t, client.Scopes, 1)
	assert.Equal(t, scopeID, client.Scopes[0].ScopeID)
	assert.Equal(t, "Protocol Resource", client.Scopes[0].ProtectedResourceName)

	missing, err := repo.FindById(ctx, uuid.NewString())
	require.NoError(t, err)
	assert.Nil(t, missing)
}

func TestPostgresSessionClientSecretRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, subscriptionID := integrationtest.BackofficeSubscription(t, pool)
	clientID, secretID := uuid.NewString(), uuid.NewString()
	_, err := pool.Exec(ctx, `INSERT INTO security.client_secrets (id,client_id,subscription_id,secret_hash) VALUES ($1,$2,$3,$4)`, secretID, clientID, subscriptionID, "argon2-hash")
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM security.client_secrets WHERE id=$1`, secretID) })

	repo := infrastructure.NewPostgresSessionClientSecretRepository(pool)
	hash, err := repo.FindHashByClientId(ctx, clientID)
	require.NoError(t, err)
	assert.Equal(t, "argon2-hash", hash)
	_, err = repo.FindHashByClientId(ctx, uuid.NewString())
	assert.Error(t, err)
}
