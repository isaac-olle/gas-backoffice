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

	t.Run("decodes the idp state, exchanges the code, and creates an already active session", func(t *testing.T) {
		t.Parallel()
		_, client, payload := newFixture(t)
		idpClient := &stubIdpClient{auth: successfulAuthentication()}
		sessionRepo := &stubSessionRepo{}

		svc := create_session.NewService(
			&stubStateCodec{payload: payload},
			domain_services.NewClientFinder(&stubClientRepo{client: client}),
			protocolSharedDS.NewIdentityProviderFinder(&stubIdpRepo{idp: aGoogleIdp()}),
			idpClient,
			sessionRepo,
			fixedCodeGen{},
			noopEventBus{},
		)

		redirect, err := svc.CreateSession(context.Background(), "the-code", "opaque-state", "google")

		require.NoError(t, err)
		assert.Equal(t, "the-code", idpClient.gotCode)
		require.NotNil(t, sessionRepo.saved)
		assert.Equal(t, "active", sessionRepo.saved.ToPrimitives().SessionState)
		assert.Equal(t, redirectURI+"?code=gas-code&state=client-state", redirect)
	})

	t.Run("rejects an invalid or expired idp state", func(t *testing.T) {
		t.Parallel()
		_, client, _ := newFixture(t)
		svc := create_session.NewService(
			&stubStateCodec{err: assert.AnError},
			domain_services.NewClientFinder(&stubClientRepo{client: client}),
			protocolSharedDS.NewIdentityProviderFinder(&stubIdpRepo{idp: aGoogleIdp()}),
			&stubIdpClient{},
			&stubSessionRepo{},
			fixedCodeGen{},
			noopEventBus{},
		)

		_, err := svc.CreateSession(context.Background(), "code", "bad-state", "google")

		assert.Error(t, err)
	})

	t.Run("rejects a state whose idp reference name does not match the callback path", func(t *testing.T) {
		t.Parallel()
		_, client, payload := newFixture(t)
		idpClient := &stubIdpClient{}
		svc := create_session.NewService(
			&stubStateCodec{payload: payload},
			domain_services.NewClientFinder(&stubClientRepo{client: client}),
			protocolSharedDS.NewIdentityProviderFinder(&stubIdpRepo{idp: aGoogleIdp()}),
			idpClient,
			&stubSessionRepo{},
			fixedCodeGen{},
			noopEventBus{},
		)

		_, err := svc.CreateSession(context.Background(), "code", "opaque-state", "some-other-idp")

		assert.Error(t, err)
		assert.Empty(t, idpClient.gotCode)
	})

	t.Run("returns error without saving when the idp exchange fails", func(t *testing.T) {
		t.Parallel()
		_, client, payload := newFixture(t)
		sessionRepo := &stubSessionRepo{}
		svc := create_session.NewService(
			&stubStateCodec{payload: payload},
			domain_services.NewClientFinder(&stubClientRepo{client: client}),
			protocolSharedDS.NewIdentityProviderFinder(&stubIdpRepo{idp: aGoogleIdp()}),
			&stubIdpClient{err: assert.AnError},
			sessionRepo,
			fixedCodeGen{},
			noopEventBus{},
		)

		_, err := svc.CreateSession(context.Background(), "code", "opaque-state", "google")

		assert.Error(t, err)
		assert.Nil(t, sessionRepo.saved)
	})

	t.Run("returns error when the client cannot be found", func(t *testing.T) {
		t.Parallel()
		_, _, payload := newFixture(t)
		svc := create_session.NewService(
			&stubStateCodec{payload: payload},
			domain_services.NewClientFinder(&stubClientRepo{client: nil}),
			protocolSharedDS.NewIdentityProviderFinder(&stubIdpRepo{idp: aGoogleIdp()}),
			&stubIdpClient{},
			&stubSessionRepo{},
			fixedCodeGen{},
			noopEventBus{},
		)

		_, err := svc.CreateSession(context.Background(), "code", "opaque-state", "google")

		assert.Error(t, err)
	})
}
