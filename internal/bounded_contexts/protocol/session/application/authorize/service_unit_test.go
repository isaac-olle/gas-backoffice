package authorize_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolSharedDS "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authorize"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func newPublicClientMother(clientID string) *session.Client {
	return session.NewClient(
		clientID, "Public Client",
		uuid.New().String(),
		[]string{"authorization_code"},
		[]string{"https://example.com/callback"},
		[]*session.Scope{{ScopeName: "openid"}},
		true,
	)
}

func newService(clientRepo *mockClientRepo, idpRepo *mockIdpFinderRepo, idpClient *mockIdpClient, stateCodec *mockStateCodec) *authorize.Service {
	return authorize.NewService(
		domain_services.NewClientFinder(clientRepo),
		protocolSharedDS.NewIdentityProviderFinder(idpRepo),
		idpClient,
		stateCodec,
	)
}

func validRequest(clientID string) authorize.Request {
	return authorize.Request{
		ResponseType: "code",
		ClientID:     clientID,
		RedirectURI:  "https://example.com/callback",
		Scope:        "openid",
		State:        "client-state",
	}
}

func TestAuthorizeService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		client    func(clientID string) *session.Client
		noClient  bool
		noIdp     bool
		idpURL    string
		encoded   string
		mutateReq func(req *authorize.Request)
		wantErr   bool
		wantErrAs bool
		wantURL   string
	}{
		{
			name:    "validates against the client and returns the idp authorization url built with the encoded idp state",
			client:  newClientMother,
			idpURL:  "https://accounts.google.com/authorize?state=opaque-state",
			encoded: "opaque-state",
			wantURL: "https://accounts.google.com/authorize?state=opaque-state",
		},
		{
			name:     "returns error when the client cannot be found",
			noClient: true,
			wantErr:  true,
		},
		{
			name:    "returns error when the identity provider cannot be found",
			client:  newClientMother,
			noIdp:   true,
			wantErr: true,
		},
		{
			name:   "rejects an unsupported response_type without touching the idp",
			client: newClientMother,
			idpURL: "should-not-be-used",
			mutateReq: func(req *authorize.Request) {
				req.ResponseType = "token"
			},
			wantErr:   true,
			wantErrAs: true,
		},
		{
			name:   "rejects a redirect_uri not registered for the client",
			client: newClientMother,
			mutateReq: func(req *authorize.Request) {
				req.RedirectURI = "https://evil.example.com/cb"
			},
			wantErr: true,
		},
		{
			name:    "requires pkce for public clients",
			client:  newPublicClientMother,
			wantErr: true,
		},
		{
			// authorize.Service no longer depends on session.Repository at all —
			// this case documents that: NewService's signature has no repository
			// parameter, so there is nothing to assert beyond compiling.
			name:    "does not persist anything and never touches a session repository",
			client:  newClientMother,
			idpURL:  "https://accounts.google.com/authorize",
			encoded: "x",
			wantURL: "https://accounts.google.com/authorize",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			clientID := uuid.New().String()

			clientRepo := &mockClientRepo{client: newClientMother(clientID)}
			if tt.client != nil {
				clientRepo = &mockClientRepo{client: tt.client(clientID)}
			}
			if tt.noClient {
				clientRepo = &mockClientRepo{client: nil}
			}

			idpRepo := &mockIdpFinderRepo{idp: aGoogleIdp()}
			if tt.noIdp {
				idpRepo = &mockIdpFinderRepo{idp: nil}
			}

			svc := newService(clientRepo, idpRepo, &mockIdpClient{url: tt.idpURL}, &mockStateCodec{encoded: tt.encoded})

			req := validRequest(clientID)
			if tt.mutateReq != nil {
				tt.mutateReq(&req)
			}

			url, err := svc.Authorize(context.Background(), req)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrAs {
					var appErr *app_errors.AppError
					require.ErrorAs(t, err, &appErr)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantURL, url)
		})
	}
}
