package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresIdentityProviderRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	repo := infrastructure.NewPostgresIdentityProviderRepository(pool)
	id := ids.NewIdentityProviderId()
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.identity_providers WHERE id=$1`, id.String()) })

	idp, err := identity_providers.NewIdentityProvider(id, "integration-"+id.String(), "Integration IdP", "https://issuer.example/"+id.String(), "client-id", "client-secret", "https://issuer.example/authorize", "https://issuer.example/token", "", "https://issuer.example/jwks", []string{"openid", "email"}, "https://app.example/callback")
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, idp))

	byID, err := repo.FindById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.Equal(t, "client-secret", byID.ToPrimitives().ClientSecret)
	assert.Empty(t, byID.ToPrimitives().UserinfoEndpoint)

	byReferenceName, err := repo.FindByReferenceName(ctx, "integration-"+id.String())
	require.NoError(t, err)
	require.NotNil(t, byReferenceName)
	assert.Equal(t, id.String(), byReferenceName.ToPrimitives().Id)

	missing, err := repo.FindById(ctx, ids.NewIdentityProviderId())
	require.NoError(t, err)
	assert.Nil(t, missing)
}
