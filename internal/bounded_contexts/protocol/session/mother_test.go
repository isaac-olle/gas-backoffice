package session_test

import (
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
)

// mockCodeGenerator stubs the code generation interface.
type mockCodeGenerator struct{ code string }

func (m *mockCodeGenerator) Generate() (string, error) { return m.code, nil }

// newSessionMother creates a valid auth-code Session with sensible defaults.
func newSessionMother(t *testing.T) *session.AuthorizationCodeSession {
	t.Helper()
	sess, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(),
		uuid.New().String(),
		"https://example.com/callback",
		"openid",
		"random-state",
		"code",
		&mockCodeGenerator{code: "auth-code-42"},
		"", "",
	)
	require.NoError(t, err)
	return sess
}

// newClientMother creates a Client with the given clientID, redirect URIs and scopes.
func newClientMother(clientID, redirectURI string, flows, scopes []string) *session.Client {
	scopeObjs := make([]*session.Scope, len(scopes))
	for i, s := range scopes {
		scopeObjs[i] = &session.Scope{ScopeName: s}
	}
	return session.NewClient(
		clientID, "Test Client",
		uuid.New().String(), uuid.New().String(),
		flows,
		[]string{redirectURI},
		scopeObjs,
	)
}
