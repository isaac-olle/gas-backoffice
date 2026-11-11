package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresScopeRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawSubscriptionID := integrationtest.BackofficeSubscription(t, pool)
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	require.NoError(t, err)
	protectedResourceID := ids.NewProtectedResourceId()
	_, err = pool.Exec(ctx, `INSERT INTO backoffice.protected_resources (id,name,subscription_id) VALUES ($1,$2,$3)`, protectedResourceID.String(), "Scope Parent", subscriptionID.String())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.protected_resources WHERE id=$1`, protectedResourceID.String())
	})

	repo := infrastructure.NewPostgresScopeRepository(pool)
	id := ids.NewScopeId()
	t.Cleanup(func() { _ = repo.HardDelete(ctx, id) })
	ownership := scopes.NewOwnershipFromVos(protectedResourceID, subscriptionID)
	scope, err := scopes.NewScope(id, "Integration Scope", ownership, false)
	require.NoError(t, err)
	clientID := ids.NewClientId()
	scope.AddAllowedClient(clientID)
	require.NoError(t, repo.Save(ctx, scope))

	loaded, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, protectedResourceID, loaded.ParentProtectedResourceId())
	assert.Equal(t, []string{clientID.String()}, loaded.ToPrimitives().AllowedClientIds)

	matched, err := repo.Matching(ctx, criteria.Criteria{Filters: []criteria.Filter{{Field: "id", Operator: criteria.Equals, Value: id.String()}}})
	require.NoError(t, err)
	require.Len(t, matched, 1)

	require.NoError(t, repo.SoftDelete(ctx, id))
	deleted, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestPostgresClientProjectionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawSubscriptionID := integrationtest.BackofficeSubscription(t, pool)
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	require.NoError(t, err)
	clientID := ids.NewClientId()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.clients WHERE id=$1`, clientID.String()) })
	_, err = pool.Exec(ctx, `INSERT INTO backoffice.clients (id,name,subscription_id) VALUES ($1,$2,$3)`, clientID.String(), "Projection Client", subscriptionID.String())
	require.NoError(t, err)

	repo := infrastructure.NewClientProjectionRepository(pool)
	projection, err := repo.FindClientProjection(ctx, clientID)
	require.NoError(t, err)
	require.NotNil(t, projection)
	assert.Equal(t, clientID, projection.ClientId)
	assert.Equal(t, subscriptionID, projection.SubscriptionId)

	missing, err := repo.FindClientProjection(ctx, ids.NewClientId())
	require.NoError(t, err)
	assert.Nil(t, missing)

	_, err = pool.Exec(ctx, `UPDATE backoffice.clients SET deleted_at=NOW() WHERE id=$1`, clientID.String())
	require.NoError(t, err)
	deleted, err := repo.FindClientProjection(ctx, clientID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
