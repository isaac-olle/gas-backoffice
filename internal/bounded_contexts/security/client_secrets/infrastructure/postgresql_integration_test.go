package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresClientSecretRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawSubscriptionID := integrationtest.BackofficeSubscription(t, pool)
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	require.NoError(t, err)
	repo := infrastructure.NewPostgresClientSecretRepository(pool)
	clientID := ids.NewClientId()
	secretID := ids.NewClientSecretId()
	t.Cleanup(func() { _ = repo.HardDelete(ctx, secretID) })

	secret, err := client_secrets.NewClientSecretForSubscription(secretID, clientID, subscriptionID, "secret-value", infrastructure.NewArgon2SecretHasher())
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, secret))

	byID, err := repo.SearchById(ctx, secretID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.NotEqual(t, "secret-value", byID.ToPrimitives().SecretHash)

	byClient, err := repo.FindByClientId(ctx, clientID)
	require.NoError(t, err)
	require.NotNil(t, byClient)
	assert.Equal(t, secretID.String(), byClient.ToPrimitives().Id)

	verifier := infrastructure.NewArgon2SecretVerifier(repo)
	require.NoError(t, verifier.Verify(ctx, clientID.String(), "secret-value"))
	assert.Error(t, verifier.Verify(ctx, clientID.String(), "wrong-secret"))

	require.NoError(t, repo.SoftDelete(ctx, secretID))
	_, err = repo.SearchById(ctx, secretID)
	assert.Error(t, err)
}
