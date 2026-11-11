package persistance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

type sessionIntegrationEncrypter struct{}

func (sessionIntegrationEncrypter) Encrypt(value string) (string, error) {
	return "encrypted:" + value, nil
}
func (sessionIntegrationEncrypter) Decrypt(value string) (string, error) {
	const prefix = "encrypted:"
	if len(value) < len(prefix) || value[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid encrypted value")
	}
	return value[len(prefix):], nil
}
func (sessionIntegrationEncrypter) Seal(value []byte) ([]byte, error) { return value, nil }
func (sessionIntegrationEncrypter) Open(value []byte) ([]byte, error) { return value, nil }

func TestPostgresAuthorizationCodeSessionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	idpID := ids.NewIdentityProviderId().String()
	sessionID := ids.NewSessionId()
	_, err := pool.Exec(ctx, `INSERT INTO backoffice.identity_providers (id,reference_name,name,issuer,client_id,authorization_endpoint,token_endpoint,jwks_uri) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, idpID, "session-"+idpID, "Integration IdP", "https://issuer.example/"+idpID, "client", "https://issuer.example/auth", "https://issuer.example/token", "https://issuer.example/jwks")
	require.NoError(t, err)
	repo := persistance.NewPostgresAuthorizationCodeSessionRepository(pool, sessionIntegrationEncrypter{})
	t.Cleanup(func() {
		_ = repo.HardDelete(ctx, sessionID)
		_, _ = pool.Exec(ctx, `DELETE FROM backoffice.identity_providers WHERE id=$1`, idpID)
	})
	refreshToken := "refresh-token"
	aggregate, err := session.AuthorizationCodeSessionFromPrimitives(&session.AuthorizationCodeSessionPrimitives{
		BaseSessionPrimitives: &session.BaseSessionPrimitives{ClientID: ids.NewClientId().String(), Scope: "openid", GrantType: "authorization_code", ProtectedResourceId: ids.NewProtectedResourceId().String()},
		SessionID:             sessionID.String(), RedirectURI: "https://app.example/callback", State: "state", ResponseTypes: []string{"code"},
		Code: "integration-code", SessionState: "active", Consumed: false,
		Identity: session.UserIdentityPrimitives{IdpId: idpID, IDToken: "id-token", AccessToken: "access-token", RefreshToken: &refreshToken, Claims: session.IdpClaimsPrimitives{Sub: "subject", Email: "user@example.test"}, ExpiresAt: time.Now().UTC().Add(time.Hour)},
	})
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, aggregate))

	byID, err := repo.FindById(ctx, sessionID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.Equal(t, "id-token", byID.ToPrimitives().Identity.IDToken)
	assert.Equal(t, "access-token", byID.ToPrimitives().Identity.AccessToken)
	assert.Equal(t, refreshToken, *byID.ToPrimitives().Identity.RefreshToken)

	byCode, err := repo.FindByCode(ctx, "integration-code")
	require.NoError(t, err)
	require.NotNil(t, byCode)
	assert.Equal(t, sessionID.String(), byCode.ToPrimitives().SessionID)

	require.NoError(t, repo.SoftDelete(ctx, sessionID))
	deleted, err := repo.FindById(ctx, sessionID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestPostgresClientCredentialsSessionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	id := ids.NewSessionId()
	sessionID, err := session.NewClientCredentialsSession(id, ids.NewClientId().String(), ids.NewUserId().String(), "read write", ids.NewProtectedResourceId().String())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM protocol.client_credentials_sessions WHERE id=$1`, id.String())
	})
	repo := persistance.NewPostgresClientCredentialsSessionRepository(pool)
	require.NoError(t, repo.Save(ctx, sessionID))

	loaded, err := repo.FindById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "read write", loaded.TokenClaims().Scope)
	assert.Equal(t, id.String(), loaded.ToPrimitives().ID)

	missing, err := repo.FindById(ctx, ids.NewSessionId())
	require.NoError(t, err)
	assert.Nil(t, missing)
}
