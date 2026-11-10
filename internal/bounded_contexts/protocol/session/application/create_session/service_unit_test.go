package create_session_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	protocolSharedDS "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// --- stubs ---

type stubStateCodec struct {
	payload session.IdpStatePayload
	err     error
}

func (s *stubStateCodec) Encode(_ session.IdpStatePayload) (string, error) { return "", nil }
func (s *stubStateCodec) Decode(_ string) (session.IdpStatePayload, error) {
	return s.payload, s.err
}

type stubClientRepo struct {
	client *session.Client
	err    error
}

func (s *stubClientRepo) FindById(_ context.Context, _ string) (*session.Client, error) {
	return s.client, s.err
}

type stubIdpRepo struct {
	idp *protocolShared.IdentityProviderConfig
	err error
}

func (s *stubIdpRepo) FindByReferenceName(_ context.Context, _ string) (*protocolShared.IdentityProviderConfig, error) {
	return s.idp, s.err
}

type stubIdpClient struct {
	gotCode string
	auth    session.IdpAuthentication
	err     error
}

func (s *stubIdpClient) ExchangeCode(_ context.Context, _ *protocolShared.IdentityProviderConfig, code string) (session.IdpAuthentication, error) {
	s.gotCode = code
	if s.err != nil {
		return session.IdpAuthentication{}, s.err
	}
	return s.auth, nil
}

type stubSessionRepo struct {
	saved *session.AuthorizationCodeSession
	err   error
}

func (s *stubSessionRepo) Save(_ context.Context, sess *session.AuthorizationCodeSession) error {
	if s.err != nil {
		return s.err
	}
	s.saved = sess
	return nil
}
func (s *stubSessionRepo) FindById(context.Context, ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (s *stubSessionRepo) FindByCode(context.Context, string) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (s *stubSessionRepo) SoftDelete(context.Context, ids.SessionId) error { return nil }
func (s *stubSessionRepo) HardDelete(context.Context, ids.SessionId) error { return nil }

type noopEventBus struct{}

func (noopEventBus) Publish(context.Context, []event.Event) {}

type fixedCodeGen struct{}

func (fixedCodeGen) Generate() (string, error) { return "gas-code", nil }

// --- fixture ---

const redirectURI = "https://client.example.com/cb"

func newFixture(t *testing.T) (clientID string, client *session.Client, payload session.IdpStatePayload) {
	t.Helper()
	clientID = uuid.New().String()
	client = session.NewClient(
		clientID, "Test Client", uuid.New().String(),
		[]string{"authorization_code"}, []string{redirectURI},
		[]*session.Scope{session.NewScope(uuid.New().String(), "openid", uuid.New().String(), "Resource")},
		false,
	)
	payload = session.IdpStatePayload{
		ClientID:         clientID,
		RedirectURI:      redirectURI,
		Scope:            "openid",
		ResponseType:     "code",
		ClientState:      "client-state",
		IdpReferenceName: "google",
	}
	return clientID, client, payload
}

func aGoogleIdp() *protocolShared.IdentityProviderConfig {
	return protocolShared.NewIdentityProviderConfig(
		uuid.New().String(), "Google", "https://accounts.google.com", "client-id", "secret",
		"", "https://oauth2.googleapis.com/token", "https://www.googleapis.com/oauth2/v3/certs",
		nil,
		"https://api.example.com/v1/oidc/callback/google",
	)
}

func successfulAuthentication() session.IdpAuthentication {
	return session.IdpAuthentication{
		Tokens: session.IdpTokens{IDToken: "id-token", AccessToken: "access-token", ExpiresAt: time.Now().Add(time.Hour)},
		Claims: sessionvos.IdpClaims{Sub: "google-sub", Email: "user@example.com"},
	}
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		stateErr     error
		noClient     bool
		idpClientErr error
		callbackIdp  string
		code         string
		wantErr      bool
		check        func(t *testing.T, idpClient *stubIdpClient, sessionRepo *stubSessionRepo, redirect string)
	}{
		{
			name:        "decodes the idp state, exchanges the code, and creates an already active session",
			callbackIdp: "google",
			code:        "the-code",
			check: func(t *testing.T, idpClient *stubIdpClient, sessionRepo *stubSessionRepo, redirect string) {
				assert.Equal(t, "the-code", idpClient.gotCode)
				require.NotNil(t, sessionRepo.saved)
				assert.Equal(t, "active", sessionRepo.saved.ToPrimitives().SessionState)
				assert.Equal(t, redirectURI+"?code=gas-code&state=client-state", redirect)
			},
		},
		{
			name:        "rejects an invalid or expired idp state",
			stateErr:    assert.AnError,
			callbackIdp: "google",
			code:        "code",
			wantErr:     true,
		},
		{
			name:        "rejects a state whose idp reference name does not match the callback path",
			callbackIdp: "some-other-idp",
			code:        "code",
			wantErr:     true,
			check: func(t *testing.T, idpClient *stubIdpClient, sessionRepo *stubSessionRepo, redirect string) {
				assert.Empty(t, idpClient.gotCode)
			},
		},
		{
			name:         "returns error without saving when the idp exchange fails",
			idpClientErr: assert.AnError,
			callbackIdp:  "google",
			code:         "code",
			wantErr:      true,
			check: func(t *testing.T, idpClient *stubIdpClient, sessionRepo *stubSessionRepo, redirect string) {
				assert.Nil(t, sessionRepo.saved)
			},
		},
		{
			name:        "returns error when the client cannot be found",
			noClient:    true,
			callbackIdp: "google",
			code:        "code",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, client, payload := newFixture(t)
			if tt.noClient {
				client = nil
			}

			idpClient := &stubIdpClient{err: tt.idpClientErr, auth: successfulAuthentication()}
			sessionRepo := &stubSessionRepo{}

			svc := create_session.NewService(
				&stubStateCodec{payload: payload, err: tt.stateErr},
				domain_services.NewClientFinder(&stubClientRepo{client: client}),
				protocolSharedDS.NewIdentityProviderFinder(&stubIdpRepo{idp: aGoogleIdp()}),
				idpClient,
				sessionRepo,
				fixedCodeGen{},
				noopEventBus{},
			)

			redirect, err := svc.CreateSession(context.Background(), tt.code, "opaque-state", tt.callbackIdp)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.check != nil {
					tt.check(t, idpClient, sessionRepo, redirect)
				}
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, idpClient, sessionRepo, redirect)
			}
		})
	}
}
