package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresProtectedResourceRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawSubscriptionID := integrationtest.BackofficeSubscription(t, pool)
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	require.NoError(t, err)
	repo := infrastructure.NewPostgresProtectedResourceRepository(pool)
	id := ids.NewProtectedResourceId()
	t.Cleanup(func() { _ = repo.HardDelete(ctx, id) })

	resource, err := protected_resources.NewProtectedResource(id, "Integration Resource", protected_resources.NewOwnership(subscriptionID), false)
	require.NoError(t, err)
	scopeID := ids.NewScopeId()
	resource.AddScopeId(scopeID)
	require.NoError(t, repo.Save(ctx, resource))

	loaded, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, subscriptionID, loaded.SubscriptionId())
	assert.Equal(t, []string{scopeID.String()}, loaded.ToPrimitives().ScopeIds)

	matched, err := repo.Matching(ctx, criteria.Criteria{Filters: []criteria.Filter{{Field: "id", Operator: criteria.Equals, Value: id.String()}}})
	require.NoError(t, err)
	require.Len(t, matched, 1)

	require.NoError(t, repo.SoftDelete(ctx, id))
	deleted, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
