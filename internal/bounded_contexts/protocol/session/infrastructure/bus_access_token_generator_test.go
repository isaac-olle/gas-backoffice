package infrastructure_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func newResolvedBaseSession(t *testing.T) *session.BaseSession {
	t.Helper()
	sess, err := session.AuthorizationCodeSessionFromPrimitives(&session.AuthorizationCodeSessionPrimitives{
		BaseSessionPrimitives: &session.BaseSessionPrimitives{
			ClientID:            uuid.New().String(),
			Scope:               "openid",
			GrantType:           "authorization_code",
			ProtectedResourceId: uuid.New().String(),
		},
		SessionID:     ids.NewSessionId().String(),
		RedirectURI:   "https://example.com/cb",
		State:         "state",
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
	return sess.BaseSession
}

func TestBusAccessTokenGenerator_GenerateAccessToken(t *testing.T) {
	t.Parallel()

	t.Run("returns error when authenticated user not resolved", func(t *testing.T) {
		t.Parallel()
		gen := infrastructure.NewBusAccessTokenGenerator(&mockQueryBus{})

		_, err := gen.GenerateAccessToken(context.Background(), newResolvedBaseSession(t), session.NewUserIdentityProjection("", ""))

		require.Error(t, err)
	})

	t.Run("dispatches to query bus when authenticated user is resolved", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{result: "signed-jwt"}
		gen := infrastructure.NewBusAccessTokenGenerator(bus)

		token, err := gen.GenerateAccessToken(context.Background(), newResolvedBaseSession(t), session.NewUserIdentityProjection(uuid.New().String(), ""))

		require.NoError(t, err)
		assert.Equal(t, "signed-jwt", token)
	})
}
