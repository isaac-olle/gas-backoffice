package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	userProjection "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/user_projection"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userInfrastructure "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresSessionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	repo := infrastructure.NewPostgresSessionRepository(pool)
	session, err := auth.NewSession(ids.NewUserId().String(), vo.SubscriptionUserRole.String(), ids.NewUserId().String())
	require.NoError(t, err)
	sessionID, err := ids.NewSessionIdFromString(session.ToPrimitives().SessionId)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.sessions WHERE id=$1`, sessionID.String()) })

	require.NoError(t, repo.SaveSession(ctx, session))
	active, err := repo.IsActive(ctx, sessionID)
	require.NoError(t, err)
	assert.True(t, active)

	require.NoError(t, repo.Logout(ctx, sessionID))
	active, err = repo.IsActive(ctx, sessionID)
	require.NoError(t, err)
	assert.False(t, active)
}

func TestPostgresIdpConfigRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	id := ids.NewIdentityProviderId().String()
	referenceName := "google"
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.identity_providers WHERE id=$1`, id) })
	_, err := pool.Exec(ctx, `INSERT INTO backoffice.identity_providers (id,reference_name,name,issuer,client_id,authorization_endpoint,token_endpoint,jwks_uri) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, id, referenceName, "Integration IdP", "https://issuer.example/"+id, "client", "https://issuer.example/auth", "https://issuer.example/token", "https://issuer.example/jwks")
	require.NoError(t, err)

	repo := infrastructure.NewPostgresIdpConfigRepository(pool)
	idpReferenceName, err := vos.NewFromString(referenceName)
	require.NoError(t, err)
	config, err := repo.FindByReferenceName(ctx, idpReferenceName)
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, referenceName, config.IdpReferenceName)

	_, err = pool.Exec(ctx, `UPDATE backoffice.identity_providers SET enabled=FALSE WHERE id=$1`, id)
	require.NoError(t, err)
	config, err = repo.FindByReferenceName(ctx, idpReferenceName)
	require.NoError(t, err)
	assert.Nil(t, config)
}

func TestPostgresUserProjectionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	repo := userProjection.NewUserProjectionRepository(pool)
	id := ids.NewUserId()
	email := "projection-" + id.String() + "@example.test"
	user, err := users.NewSubscriptionUser(id, "Projection User", email)
	require.NoError(t, err)
	adminID := ids.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminID.String(), adminID.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	ctx = vo.WithAuthor(ctx, author)
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.users WHERE id=$1`, id.String()) })
	require.NoError(t, userInfrastructure.NewPostgresUserRepository(pool).Save(ctx, user))

	projection, err := repo.FindByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, projection)
	assert.Equal(t, id.String(), projection.UserId)
	assert.Equal(t, adminID.String(), projection.AdminId)
	assert.Equal(t, vo.SubscriptionUserRole.String(), projection.SystemRole)

	_, err = repo.FindByEmail(ctx, "missing-"+email)
	assert.Error(t, err)
}
