package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresClientProjectionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, subscriptionID := integrationtest.BackofficeSubscription(t, pool)
	clientID := ids.NewClientId()
	_, err := pool.Exec(ctx, `INSERT INTO backoffice.clients (id,name,subscription_id) VALUES ($1,$2,$3)`, clientID.String(), "Projection Client", subscriptionID)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.clients WHERE id=$1`, clientID.String()) })

	projection, err := infrastructure.NewClientProjectionRepository(pool).Find(ctx, clientID)
	require.NoError(t, err)
	require.NotNil(t, projection)
	assert.Equal(t, clientID.String(), projection.ClientId)
	assert.Equal(t, subscriptionID, projection.SubscriptionId)

	_, err = infrastructure.NewClientProjectionRepository(pool).Find(ctx, ids.NewClientId())
	assert.Error(t, err)
}
