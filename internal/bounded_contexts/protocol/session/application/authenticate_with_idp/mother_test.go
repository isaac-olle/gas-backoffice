package authenticate_with_idp_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// ── session repo mock (for SessionFinder) ──────────────────────────────────────

type mockSessionRepo struct {
	session *session.AuthorizationCodeSession
	err     error
}

func (m *mockSessionRepo) Save(context.Context, *session.AuthorizationCodeSession) error { return nil }
func (m *mockSessionRepo) FindById(context.Context, ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return m.session, m.err
}
func (m *mockSessionRepo) FindByCode(context.Context, string) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (m *mockSessionRepo) SoftDelete(context.Context, ids.SessionId) error { return nil }
func (m *mockSessionRepo) HardDelete(context.Context, ids.SessionId) error { return nil }

// ── idp repo mock (for IdentityProviderFinder) ─────────────────────────────────

type mockIdpRepo struct {
	idp *session.IdentityProvider
	err error
}

func (m *mockIdpRepo) FindById(context.Context, string) (*session.IdentityProvider, error) {
	return m.idp, m.err
}

// ── idp client mock ────────────────────────────────────────────────────────────

type mockIdpClient struct {
	url string
	err error
}

func (m *mockIdpClient) BuildAuthorizationURL(*session.IdentityProvider, *session.AuthorizationCodeSession) (string, error) {
	return m.url, m.err
}

// ── fixtures ────────────────────────────────────────────────────────────────────

func aSession(t *testing.T) *session.AuthorizationCodeSession {
	t.Helper()
	sess, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(), uuid.New().String(), "https://fe.example.com/cb", "openid", "state", "code", stubCodeGen{}, "", "",
	)
	require.NoError(t, err)
	return sess
}

func aGoogleIDP() *session.IdentityProvider {
	return session.NewIdentityProvider(
		"idp-123", "google", "https://accounts.google.com", "client-abc", "secret",
		"https://accounts.google.com/o/oauth2/v2/auth", "https://oauth2.googleapis.com/token",
		"", "https://www.googleapis.com/oauth2/v3/certs", []string{"openid"},
	)
}

type stubCodeGen struct{}

func (stubCodeGen) Generate() (string, error) { return "code", nil }

var errBoom = errors.New("boom")
