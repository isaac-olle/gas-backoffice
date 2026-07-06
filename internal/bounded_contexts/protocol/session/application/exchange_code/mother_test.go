package exchange_code_test

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type mockSessionRepo struct {
	byCode  *session.AuthorizationCodeSession
	byId    *session.AuthorizationCodeSession
	saveErr error
}

func (m *mockSessionRepo) Save(_ context.Context, _ *session.AuthorizationCodeSession) error {
	return m.saveErr
}
func (m *mockSessionRepo) FindById(_ context.Context, _ ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return m.byId, nil
}
func (m *mockSessionRepo) FindByCode(_ context.Context, _ string) (*session.AuthorizationCodeSession, error) {
	return m.byCode, nil
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

type mockTokenGen struct {
	token string
	err   error
}

func (m *mockTokenGen) GenerateFromSession(_ context.Context, _ *session.AuthorizationCodeSession, _ *session.Client, _ string) (string, error) {
	return m.token, m.err
}
func (m *mockTokenGen) GenerateFromClientCredentials(_ context.Context, _ *session.ClientCredentialsSession, _ string) (string, error) {
	return m.token, m.err
}

type codeGen struct{}

func (c *codeGen) Generate() (string, error) { return "auth-code", nil }

func newPKCESessionMother(t *testing.T, clientID string) *session.AuthorizationCodeSession {
	t.Helper()
	s, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(), clientID,
		"https://example.com/cb", "openid", "state",
		"code", &codeGen{}, "my-challenge", "plain",
	)
	require.NoError(t, err)
	return s
}

func newPKCEClientMother(clientID, redirectURI string) *session.Client {
	return session.NewClient(
		clientID, "TC",
		uuid.New().String(), uuid.New().String(),
		[]string{"pkce"}, []string{redirectURI},
		[]*session.Scope{
			session.NewScope(uuid.New().String(), "openid", uuid.New().String(), "Resource A"),
		},
	)
}
