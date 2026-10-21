package token_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type mockSessionRepo struct {
	byCode  *session.AuthorizationCodeSession
	findErr error
	saved   *session.AuthorizationCodeSession
	saveErr error
}

func (m *mockSessionRepo) Save(_ context.Context, s *session.AuthorizationCodeSession) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = s
	return nil
}
func (m *mockSessionRepo) FindById(_ context.Context, _ ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (m *mockSessionRepo) FindByCode(_ context.Context, _ string) (*session.AuthorizationCodeSession, error) {
	return m.byCode, m.findErr
}
func (m *mockSessionRepo) SoftDelete(_ context.Context, _ ids.SessionId) error { return nil }
func (m *mockSessionRepo) HardDelete(_ context.Context, _ ids.SessionId) error { return nil }

type mockClientRepo struct {
	client *session.Client
	err    error
}

func (m *mockClientRepo) FindById(_ context.Context, _ string) (*session.Client, error) {
	return m.client, m.err
}

type mockEventBus struct{}

func (m *mockEventBus) Publish(_ context.Context, _ []event.Event) {}

type mockAccessTokenGenerator struct {
	token string
	err   error
}

func (m *mockAccessTokenGenerator) GenerateAccessToken(_ context.Context, _ *session.BaseSession, _ *session.UserIdentityProjection) (string, error) {
	return m.token, m.err
}

type mockIDTokenGenerator struct {
	token string
	err   error
}

func (m *mockIDTokenGenerator) GenerateIDToken(_ context.Context, _ *session.BaseSession, _ *session.UserIdentityProjection) (string, error) {
	return m.token, m.err
}

type mockCCRepo struct {
	sess    *session.ClientCredentialsSession
	findErr error
	saved   *session.ClientCredentialsSession
	saveErr error
}

func (m *mockCCRepo) Save(_ context.Context, s *session.ClientCredentialsSession) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = s
	return nil
}
func (m *mockCCRepo) FindById(_ context.Context, _ ids.SessionId) (*session.ClientCredentialsSession, error) {
	if m.saved != nil {
		return m.saved, m.findErr
	}
	return m.sess, m.findErr
}

func newAuthCodeClientMother(clientID string) *session.Client {
	return session.NewClient(
		clientID, "Test Client", uuid.New().String(),
		[]string{"authorization_code"},
		[]string{"https://example.com/cb"},
		[]*session.Scope{{ScopeName: "openid"}},
		false,
	)
}

func newCCClientMother(clientID string) *session.Client {
	return session.NewClient(
		clientID, "CC Client", uuid.New().String(),
		[]string{"client_credentials"},
		nil,
		[]*session.Scope{session.NewScope("scope-1", "openid", "pr-001", "PR")},
		false,
	)
}

func newAuthSessionMother(t *testing.T, clientID, redirectURI, code string) *session.AuthorizationCodeSession {
	t.Helper()
	sess, err := session.AuthorizationCodeSessionFromPrimitives(&session.AuthorizationCodeSessionPrimitives{
		BaseSessionPrimitives: &session.BaseSessionPrimitives{
			ClientID:  clientID,
			Scope:     "openid",
			GrantType: "authorization_code",
		},
		SessionID:     ids.NewSessionId().String(),
		RedirectURI:   redirectURI,
		State:         "state",
		ResponseTypes: []string{"code"},
		Code:          code,
		SessionState:  "authorizing",
		Identity: session.UserIdentityPrimitives{
			IdpId:       ids.NewIdentityProviderId().String(),
			IDToken:     "id-token",
			AccessToken: "access-token",
			Claims:      session.IdpClaimsPrimitives{Sub: "google-sub", Email: "user@example.com"},
			ExpiresAt:   time.Now().Add(time.Hour),
		},
	})
	require.NoError(t, err)
	return sess
}

type fakeUnitOfWork struct{}

func (fakeUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
