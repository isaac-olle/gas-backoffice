package build_authorization_response_test

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

// mockSessionRepo stubs session.Repository for build_authorization_response tests.
type mockSessionRepo struct {
	byId *session.AuthorizationCodeSession
	err  error
}

func (m *mockSessionRepo) Save(_ context.Context, _ *session.AuthorizationCodeSession) error {
	return nil
}
func (m *mockSessionRepo) FindById(_ context.Context, _ ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return m.byId, m.err
}
func (m *mockSessionRepo) FindByCode(_ context.Context, _ string) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (m *mockSessionRepo) SoftDelete(_ context.Context, _ ids.SessionId) error { return nil }
func (m *mockSessionRepo) HardDelete(_ context.Context, _ ids.SessionId) error { return nil }

type codeGen struct{}

func (g *codeGen) Generate() (string, error) { return "bar-code", nil }

// newSessionMother builds a valid session for testing.
func newSessionMother(t *testing.T) *session.AuthorizationCodeSession {
	t.Helper()
	s, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(), uuid.New().String(),
		"https://example.com/cb", "openid", "mystate",
		"code", &codeGen{}, "", "",
	)
	require.NoError(t, err)
	return s
}
