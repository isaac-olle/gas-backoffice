package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	pginfra "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresSubscriptionRepositorySaveUsesAmbientTransaction(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	repo := infrastructure.NewPostgresRepository(pool)

	tests := []struct {
		name      string
		rollback  bool
		wantFound bool
	}{
		{name: "commits with the ambient transaction", wantFound: true},
		{name: "rolls back with the ambient transaction", rollback: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userID := ids.NewUserId()
			subscriptionID := ids.NewSubscriptionId()
			t.Cleanup(func() {
				_, _ = pool.Exec(ctx, `DELETE FROM backoffice.subscriptions WHERE id=$1`, subscriptionID.String())
				_, _ = pool.Exec(ctx, `DELETE FROM backoffice.users WHERE id=$1`, userID.String())
			})

			tx, err := pool.Begin(ctx)
			require.NoError(t, err)
			defer func() { _ = tx.Rollback(ctx) }()
			_, err = tx.Exec(ctx, `INSERT INTO backoffice.users (id,name,email) VALUES ($1,$2,$3)`, userID.String(), "Integration Owner", "owner-"+userID.String()+"@example.test")
			require.NoError(t, err)
			aggregate, err := subscriptions.NewSubscription(subscriptionID, "Integration Subscription", userID, false)
			require.NoError(t, err)
			transactionContext := pginfra.ContextWithTx(ctx, tx)
			require.NoError(t, repo.Save(transactionContext, aggregate))

			if tt.rollback {
				require.NoError(t, tx.Rollback(ctx))
			} else {
				require.NoError(t, tx.Commit(ctx))
			}
			loaded, err := repo.SearchById(ctx, subscriptionID)
			require.NoError(t, err)
			if tt.wantFound {
				require.NotNil(t, loaded)
				assert.Equal(t, "Integration Subscription", loaded.ToPrimitives().Name)
				assert.Equal(t, userID.String(), loaded.ToPrimitives().OwnerId)
			} else {
				assert.Nil(t, loaded)
			}
		})
	}
}

func TestPostgresSubscriptionRepositoryList(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	ownerID := ids.NewUserId()
	subscriptionID := ids.NewSubscriptionId()
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.subscriptions WHERE id=$1`, subscriptionID.String())
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.users WHERE id=$1`, ownerID.String())
	})
	_, err := pool.Exec(ctx, `INSERT INTO backoffice.users (id,name,email) VALUES ($1,$2,$3)`, ownerID.String(), "List Owner", "owner-"+ownerID.String()+"@example.test")
	require.NoError(t, err)
	repo := infrastructure.NewPostgresRepository(pool)
	aggregate, err := subscriptions.NewSubscription(subscriptionID, "Listed Subscription", ownerID, false)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, aggregate))

	listed, err := repo.List(ctx, []ids.SubscriptionId{subscriptionID}, false)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	assert.Equal(t, subscriptionID.String(), listed[0].Id)

	empty, err := repo.List(ctx, nil, false)
	require.NoError(t, err)
	assert.Empty(t, empty)
}
