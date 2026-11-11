package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresUserRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	repo := infrastructure.NewPostgresUserRepository(pool)
	id := ids.NewUserId()
	t.Cleanup(func() { _ = repo.HardDelete(ctx, id) })

	user, err := users.NewSubscriptionUser(id, "Integration User", "integration-"+id.String()+"@example.test")
	require.NoError(t, err)
	subscriptionID := ids.NewSubscriptionId()
	user.AddMembership(vo.NewSubscriptionMembership(subscriptionID, vo.UserSubscriptionRole))
	author, err := vo.NewUserAuthorFromString(ids.NewUserId().String(), ids.NewUserId().String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	ctx = vo.WithAuthor(ctx, author)

	require.NoError(t, repo.Save(ctx, user))
	loaded, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "Integration User", loaded.ToPrimitives().Name)
	assert.Equal(t, author.Id.String(), loaded.ToPrimitives().AuthorId)
	assert.Equal(t, []string{subscriptionID.String()}, loaded.ToPrimitives().SubscriptionIds)
	assert.Equal(t, vo.UserSubscriptionRole, loaded.ToPrimitives().Memberships[0].SubscriptionRole)

	matched, err := repo.Matching(ctx, criteria.Criteria{Filters: []criteria.Filter{{Field: "id", Operator: criteria.Equals, Value: id.String()}}})
	require.NoError(t, err)
	require.Len(t, matched, 1)

	stale, err := users.FromPrimitives(loaded.ToPrimitives())
	require.NoError(t, err)
	loaded, err = repo.SearchById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	name, err := vos.NewUserName("Updated Integration User")
	require.NoError(t, err)
	require.NoError(t, loaded.ChangeName(name))
	require.NoError(t, repo.Save(ctx, loaded))
	staleName, err := vos.NewUserName("Stale Integration User")
	require.NoError(t, err)
	require.NoError(t, stale.ChangeName(staleName))
	assert.ErrorIs(t, repo.Save(ctx, stale), persistance.ErrConcurrentModification)

	require.NoError(t, repo.SoftDelete(ctx, id))
	deleted, err := repo.SearchById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
