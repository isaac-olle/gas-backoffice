package persistance_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type spySessionRepo struct {
	sess          *session.AuthorizationCodeSession
	findByCodeHit int
}

func (s *spySessionRepo) Save(_ context.Context, _ *session.AuthorizationCodeSession) error {
	return nil
}
func (s *spySessionRepo) FindById(_ context.Context, _ ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (s *spySessionRepo) FindByCode(_ context.Context, _ string) (*session.AuthorizationCodeSession, error) {
	s.findByCodeHit++
	return s.sess, nil
}
func (s *spySessionRepo) SoftDelete(_ context.Context, _ ids.SessionId) error { return nil }
func (s *spySessionRepo) HardDelete(_ context.Context, _ ids.SessionId) error { return nil }

func aSession(t *testing.T) *session.AuthorizationCodeSession {
	t.Helper()
	sess, err := session.AuthorizationCodeSessionFromPrimitives(&session.AuthorizationCodeSessionPrimitives{
		BaseSessionPrimitives: &session.BaseSessionPrimitives{
			ClientID:  uuid.New().String(),
			Scope:     "openid",
			GrantType: "authorization_code",
		},
		SessionID:     ids.NewSessionId().String(),
		RedirectURI:   "https://example.com/cb",
		State:         "the-state",
		ResponseTypes: []string{"code"},
		Code:          "code",
		SessionState:  "active",
		Identity: session.UserIdentityPrimitives{
			IdpId:       ids.NewIdentityProviderId().String(),
			IDToken:     "id-token",
			AccessToken: "access-token",
			Claims:      session.IdpClaimsPrimitives{Sub: "google-sub", Email: "user@example.com"},
		},
	})
	require.NoError(t, err)
	return sess
}

func TestRequestCachedSessionRepository_FindByCode(t *testing.T) {
	t.Parallel()

	t.Run("hits the inner repository only once per cached context", func(t *testing.T) {
		t.Parallel()
		spy := &spySessionRepo{sess: aSession(t)}
		repo := persistance.NewRequestCachedSessionRepository(spy)
		ctx := persistance.WithCache(context.Background())

		first, err := repo.FindByCode(ctx, "code")
		require.NoError(t, err)
		second, err := repo.FindByCode(ctx, "code")
		require.NoError(t, err)

		assert.Same(t, first, second)
		assert.Equal(t, 1, spy.findByCodeHit)
	})

	t.Run("without a cached context, always calls the inner repository", func(t *testing.T) {
		t.Parallel()
		spy := &spySessionRepo{sess: aSession(t)}
		repo := persistance.NewRequestCachedSessionRepository(spy)

		_, err := repo.FindByCode(context.Background(), "code")
		require.NoError(t, err)
		_, err = repo.FindByCode(context.Background(), "code")
		require.NoError(t, err)

		assert.Equal(t, 2, spy.findByCodeHit)
	})

	t.Run("different contexts do not share a cache", func(t *testing.T) {
		t.Parallel()
		spy := &spySessionRepo{sess: aSession(t)}
		repo := persistance.NewRequestCachedSessionRepository(spy)

		_, err := repo.FindByCode(persistance.WithCache(context.Background()), "code")
		require.NoError(t, err)
		_, err = repo.FindByCode(persistance.WithCache(context.Background()), "code")
		require.NoError(t, err)

		assert.Equal(t, 2, spy.findByCodeHit)
	})
}
