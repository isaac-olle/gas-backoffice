package token_test

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token/flows"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func newAuthorizationCodeGrantHandler(sessionRepo *mockSessionRepo, clientRepo *mockClientRepo, accessGen *mockAccessTokenGenerator, idGen *mockIDTokenGenerator) *flows.AuthorizationCodeGrantHandler {
	return flows.NewAuthorizationCodeGrantHandler(
		domain_services.NewAuthorizationCodeSessionFinder(sessionRepo),
		domain_services.NewClientFinder(clientRepo),
		sessionRepo,
		&mockEventBus{},
		domain_services.NewAuthorizationCodeTokenResponseGenerator(accessGen, idGen),
		fakeUnitOfWork{},
	)
}

func newClientCredentialsGrantHandler(clientRepo *mockClientRepo, ccRepo *mockCCRepo, accessGen *mockAccessTokenGenerator) *flows.ClientCredentialsGrantHandler {
	return flows.NewClientCredentialsGrantHandler(
		domain_services.NewClientFinder(clientRepo),
		ccRepo,
		&mockEventBus{},
		domain_services.NewClientCredentialsSessionFinder(ccRepo),
		domain_services.NewClientCredentialsTokenResponseGenerator(accessGen),
		fakeUnitOfWork{},
	)
}

func TestAuthorizationCodeGrantHandler(t *testing.T) {
	t.Parallel()

	t.Run("consumes the code and returns the generated token", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		sess := newAuthSessionMother(t, clientID, "https://example.com/cb", "authcode")

		h := newAuthorizationCodeGrantHandler(
			&mockSessionRepo{byCode: sess},
			&mockClientRepo{client: newAuthCodeClientMother(clientID)},
			&mockAccessTokenGenerator{token: "signed.access.token"},
			&mockIDTokenGenerator{token: "signed.id.token"},
		)

		assert.True(t, h.Supports(flows.AuthorizationCodeGrantType))

		resp, err := h.Execute(context.Background(), token.Request{
			ClientID:    clientID,
			Code:        "authcode",
			RedirectURI: "https://example.com/cb",
		})

		require.NoError(t, err)
		assert.Equal(t, "signed.access.token", resp.AccessToken)
		assert.Equal(t, "signed.id.token", resp.IDToken)
	})

	t.Run("returns error when the code cannot be consumed", func(t *testing.T) {
		t.Parallel()
		h := newAuthorizationCodeGrantHandler(
			&mockSessionRepo{findErr: assert.AnError},
			&mockClientRepo{},
			&mockAccessTokenGenerator{},
			&mockIDTokenGenerator{},
		)

		_, err := h.Execute(context.Background(), token.Request{Code: "bad-code"})

		assert.Error(t, err)
	})
}

func TestClientCredentialsGrantHandler(t *testing.T) {
	t.Parallel()

	t.Run("creates the session and returns the generated token", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		h := newClientCredentialsGrantHandler(
			&mockClientRepo{client: newCCClientMother(clientID)},
			&mockCCRepo{},
			&mockAccessTokenGenerator{token: "signed.jwt.token"},
		)

		assert.True(t, h.Supports(flows.ClientCredentialsGrantType))

		resp, err := h.Execute(context.Background(), token.Request{ClientID: clientID, Scope: "openid"})

		require.NoError(t, err)
		assert.Equal(t, "signed.jwt.token", resp.AccessToken)
		assert.Empty(t, resp.IDToken)
	})

	t.Run("rejects public clients", func(t *testing.T) {
		t.Parallel()
		h := newClientCredentialsGrantHandler(&mockClientRepo{}, &mockCCRepo{}, &mockAccessTokenGenerator{})

		_, err := h.Execute(context.Background(), token.Request{IsPublic: true})

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, session.ErrUnsupportedGrantType, appErr.Code)
	})
}

func TestService(t *testing.T) {
	t.Parallel()

	t.Run("routes to the handler that supports the grant type", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		h := newClientCredentialsGrantHandler(
			&mockClientRepo{client: newCCClientMother(clientID)},
			&mockCCRepo{},
			&mockAccessTokenGenerator{token: "signed.jwt.token"},
		)
		svc := token.NewService(h)

		resp, err := svc.Execute(context.Background(), flows.ClientCredentialsGrantType, token.Request{ClientID: clientID, Scope: "openid"})

		require.NoError(t, err)
		assert.Equal(t, "signed.jwt.token", resp.AccessToken)
	})

	t.Run("returns error for an unsupported grant type", func(t *testing.T) {
		t.Parallel()
		svc := token.NewService()

		_, err := svc.Execute(context.Background(), "refresh_token", token.Request{})

		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, session.ErrUnsupportedGrantType, appErr.Code)
	})
}
