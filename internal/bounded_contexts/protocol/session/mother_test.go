package session_test

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
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
