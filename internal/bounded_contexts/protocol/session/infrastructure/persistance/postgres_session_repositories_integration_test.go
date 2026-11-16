//go:build integration

package persistance_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedpersistence "github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	sharedinfra "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
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

func TestIntegrationPostgresAuthorizationCodeSessionRepository(t *testing.T) {
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
		ExpiresAt:             time.Now().UTC().Add(15 * time.Minute),
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

	// Both exchanges read the same version before competing to persist it.
	uow := sharedinfra.NewPostgresUnitOfWork(pool)
	start := make(chan struct{})
	ready := make(chan struct{}, 2)
	update := make(chan struct{})
	results := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			results <- uow.Do(ctx, func(txCtx context.Context) error {
				loaded, loadErr := repo.FindByCode(txCtx, "integration-code")
				ready <- struct{}{}
				<-update
				if loadErr != nil {
					return loadErr
				}
				if loaded == nil {
					return errors.New("code not found")
				}
				if checkErr := loaded.CheckExchangeParams(aggregate.ClientID().String(), "https://app.example/callback"); checkErr != nil {
					return checkErr
				}
				loaded.Consume()
				return repo.Save(txCtx, loaded)
			})
		}()
	}
	close(start)
	<-ready
	<-ready
	close(update)
	successes := 0
	for range 2 {
		if exchangeErr := <-results; exchangeErr == nil {
			successes++
		} else {
			require.ErrorIs(t, exchangeErr, sharedpersistence.ErrConcurrentModification)
		}
	}
	require.Equal(t, 1, successes)
	persisted, err := repo.FindById(ctx, sessionID)
	require.NoError(t, err)
	require.NotNil(t, persisted)
	assert.EqualValues(t, 2, persisted.ToPrimitives().OccVersion)
	assert.True(t, persisted.IsConsumed())
	require.ErrorIs(t, repo.Save(ctx, byCode), sharedpersistence.ErrConcurrentModification)

	_, err = pool.Exec(ctx, `UPDATE protocol.authorization_code_sessions SET expires_at = NOW() - interval '1 second' WHERE session_id = $1`, sessionID.String())
	require.NoError(t, err)
	expired, err := repo.FindByCode(ctx, "integration-code")
	require.NoError(t, err)
	require.NotNil(t, expired)
	var expiryError *app_errors.AppError
	require.ErrorAs(t, expired.CheckExchangeParams(aggregate.ClientID().String(), "https://app.example/callback"), &expiryError)
	assert.Equal(t, session.ErrAuthorizationCodeExpired, expiryError.Code)

	require.NoError(t, repo.SoftDelete(ctx, sessionID))
	deleted, err := repo.FindById(ctx, sessionID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestIntegrationPostgresClientCredentialsSessionRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	id := ids.NewSessionId()
	sessionID, err := session.NewClientCredentialsSession(id, ids.NewClientId().String(), "read write", ids.NewProtectedResourceId().String())
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
	assert.Empty(t, loaded.ToPrimitives().UserId)

	missing, err := repo.FindById(ctx, ids.NewSessionId())
	require.NoError(t, err)
	assert.Nil(t, missing)
}
