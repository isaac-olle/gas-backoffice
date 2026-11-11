package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresClientRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, subscriptionID := integrationtest.BackofficeSubscription(t, pool)
	repo := infrastructure.NewPostgresClientRepository(pool)
	id := ids.NewClientId()
	t.Cleanup(func() { _ = repo.HardDelete(ctx, id) })

	client, err := clients.NewClient(id, "Integration Client", clients.NewOwnership(mustSubscriptionID(t, subscriptionID)), []string{"https://app.example/callback"}, []string{"authorization_code"}, false, false)
	require.NoError(t, err)
	client.AddAllowedScope(ids.NewScopeId())
	require.NoError(t, repo.Save(ctx, client))

	loaded, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, []string{"https://app.example/callback"}, loaded.ToPrimitives().RedirectUris)
	require.Len(t, loaded.ToPrimitives().AllowedScopes, 1)

	matched, err := repo.Matching(ctx, criteria.Criteria{Filters: []criteria.Filter{{Field: "id", Operator: criteria.Equals, Value: id.String()}}})
	require.NoError(t, err)
	require.Len(t, matched, 1)

	require.NoError(t, repo.SoftDelete(ctx, id))
	deleted, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

func mustSubscriptionID(t *testing.T, raw string) ids.SubscriptionId {
	t.Helper()
	id, err := ids.NewSubscriptionIdFromString(raw)
	require.NoError(t, err)
	return id
}
