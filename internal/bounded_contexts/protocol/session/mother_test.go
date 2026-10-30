package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// newUserIdentityMother builds a valid ExternalOidcUserIdentity VO with sensible defaults.
func newUserIdentityMother(t *testing.T) sessionvos.ExternalOidcUserIdentity {
	t.Helper()
	identity, err := sessionvos.NewExternalOidcUserIdentity(
		ids.NewIdentityProviderId().String(),
		"id-token",
		"access-token",
		nil,
		sessionvos.IdpClaims{Sub: "google-sub", Email: "user@example.com"},
		time.Now().Add(time.Hour),
	)
	require.NoError(t, err)
	return identity
}

// mockCodeGenerator stubs the code generation interface.
type mockCodeGenerator struct{ code string }

func (m *mockCodeGenerator) Generate() (string, error) { return m.code, nil }

// newTestClient creates a Client allowing authorization_code, the given redirectURI, and scopes.
func newTestClient(clientID, redirectURI string, scopeNames []string) *session.Client {
	prId := uuid.New().String()
	scopes := make([]*session.Scope, len(scopeNames))
	for i, s := range scopeNames {
		scopes[i] = session.NewScope(uuid.New().String(), s, prId, "Resource")
	}
	return session.NewClient(
		clientID, "Test Client",
		uuid.New().String(),
		[]string{"authorization_code"},
		[]string{redirectURI},
		scopes,
		false,
	)
}

// newSessionMother creates a valid auth-code Session with sensible defaults.
func newSessionMother(t *testing.T) *session.AuthorizationCodeSession {
	t.Helper()
	clientID := uuid.New().String()
	redirectURI := "https://example.com/callback"
	client := newTestClient(clientID, redirectURI, []string{"openid"})
	sess, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(),
		clientID,
		redirectURI,
		"openid",
		"random-state",
		"code",
		&mockCodeGenerator{code: "auth-code-42"},
		"", "",
		client,
		newUserIdentityMother(t),
	)
	require.NoError(t, err)
	return sess
}
