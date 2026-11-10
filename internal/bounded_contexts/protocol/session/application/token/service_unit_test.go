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

	tests := []struct {
		name    string
		setup   func(t *testing.T, clientID string) (*flows.AuthorizationCodeGrantHandler, token.Request)
		wantErr bool
		check   func(t *testing.T, resp *session.TokenResponse)
	}{
		{
			name: "consumes the code and returns the generated token",
			setup: func(t *testing.T, clientID string) (*flows.AuthorizationCodeGrantHandler, token.Request) {
				sess := newAuthSessionMother(t, clientID, "https://example.com/cb", "authcode")
				h := newAuthorizationCodeGrantHandler(
					&mockSessionRepo{byCode: sess},
					&mockClientRepo{client: newAuthCodeClientMother(clientID)},
					&mockAccessTokenGenerator{token: "signed.access.token"},
					&mockIDTokenGenerator{token: "signed.id.token"},
				)
				assert.True(t, h.Supports(flows.AuthorizationCodeGrantType))
				return h, token.Request{ClientID: clientID, Code: "authcode", RedirectURI: "https://example.com/cb"}
			},
			check: func(t *testing.T, resp *session.TokenResponse) {
				assert.Equal(t, "signed.access.token", resp.AccessToken)
				assert.Equal(t, "signed.id.token", resp.IDToken)
			},
		},
		{
			name: "returns error when the code cannot be consumed",
			setup: func(t *testing.T, clientID string) (*flows.AuthorizationCodeGrantHandler, token.Request) {
				h := newAuthorizationCodeGrantHandler(
					&mockSessionRepo{findErr: assert.AnError},
					&mockClientRepo{},
					&mockAccessTokenGenerator{},
					&mockIDTokenGenerator{},
				)
				return h, token.Request{Code: "bad-code"}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h, req := tt.setup(t, uuid.New().String())

			resp, err := h.Execute(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestClientCredentialsGrantHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T, clientID string) (*flows.ClientCredentialsGrantHandler, token.Request)
		wantErr     bool
		wantErrCode app_errors.ErrorCode
		check       func(t *testing.T, resp *session.TokenResponse)
	}{
		{
			name: "creates the session and returns the generated token",
			setup: func(t *testing.T, clientID string) (*flows.ClientCredentialsGrantHandler, token.Request) {
				h := newClientCredentialsGrantHandler(
					&mockClientRepo{client: newCCClientMother(clientID)},
					&mockCCRepo{},
					&mockAccessTokenGenerator{token: "signed.jwt.token"},
				)
				assert.True(t, h.Supports(flows.ClientCredentialsGrantType))
				return h, token.Request{ClientID: clientID, Scope: "openid"}
			},
			check: func(t *testing.T, resp *session.TokenResponse) {
				assert.Equal(t, "signed.jwt.token", resp.AccessToken)
				assert.Empty(t, resp.IDToken)
			},
		},
		{
			name: "rejects public clients",
			setup: func(t *testing.T, clientID string) (*flows.ClientCredentialsGrantHandler, token.Request) {
				h := newClientCredentialsGrantHandler(&mockClientRepo{}, &mockCCRepo{}, &mockAccessTokenGenerator{})
				return h, token.Request{IsPublic: true}
			},
			wantErr:     true,
			wantErrCode: session.ErrUnsupportedGrantType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h, req := tt.setup(t, uuid.New().String())

			resp, err := h.Execute(context.Background(), req)

			if tt.wantErr {
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tt.wantErrCode, appErr.Code)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T, clientID string) (*token.Service, string, token.Request)
		wantErr     bool
		wantErrCode app_errors.ErrorCode
		wantAccess  string
	}{
		{
			name: "routes to the handler that supports the grant type",
			setup: func(t *testing.T, clientID string) (*token.Service, string, token.Request) {
				h := newClientCredentialsGrantHandler(
					&mockClientRepo{client: newCCClientMother(clientID)},
					&mockCCRepo{},
					&mockAccessTokenGenerator{token: "signed.jwt.token"},
				)
				return token.NewService(h), flows.ClientCredentialsGrantType, token.Request{ClientID: clientID, Scope: "openid"}
			},
			wantAccess: "signed.jwt.token",
		},
		{
			name: "returns error for an unsupported grant type",
			setup: func(t *testing.T, clientID string) (*token.Service, string, token.Request) {
				return token.NewService(), "refresh_token", token.Request{}
			},
			wantErr:     true,
			wantErrCode: session.ErrUnsupportedGrantType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, grantType, req := tt.setup(t, uuid.New().String())

			resp, err := svc.Execute(context.Background(), grantType, req)

			if tt.wantErr {
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tt.wantErrCode, appErr.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAccess, resp.AccessToken)
		})
	}
}
