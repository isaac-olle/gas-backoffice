package infrastructure_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type codeGen struct{}

func (codeGen) Generate() (string, error) { return "the-code", nil }

func aGoogleIDP() *session.IdentityProvider {
	return session.NewIdentityProvider(
		"idp-123", "google", "https://accounts.google.com", "client-abc", "secret-xyz",
		"https://accounts.google.com/o/oauth2/v2/auth", "https://oauth2.googleapis.com/token",
		"", "https://www.googleapis.com/oauth2/v3/certs", []string{"openid", "email"},
	)
}

func aSessionWithState(t *testing.T, state string) *session.AuthorizationCodeSession {
	t.Helper()
	sess, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(), uuid.New().String(), "https://fe.example.com/cb", "openid", state, "code", codeGen{}, "", "",
	)
	require.NoError(t, err)
	return sess
}

func TestBuildAuthorizationURL(t *testing.T) {
	t.Parallel()

	client := infrastructure.NewGoogleIdentityProviderClient("https://api.example.com")
	raw, err := client.BuildAuthorizationURL(aGoogleIDP(), aSessionWithState(t, "the-state"))
	require.NoError(t, err)

	u, err := url.Parse(raw)
	require.NoError(t, err)
	assert.Equal(t, "accounts.google.com", u.Host)

	q := u.Query()
	assert.Equal(t, "code", q.Get("response_type"))
	assert.Equal(t, "client-abc", q.Get("client_id"))
	assert.Equal(t, "https://api.example.com/v1/auth/callback/idp-123", q.Get("redirect_uri"))
	assert.Equal(t, "openid email", q.Get("scope"))
	assert.Equal(t, "the-state", q.Get("state"))
	assert.NotContains(t, raw, "secret-xyz", "client secret must never appear in the authorization URL")
}

func TestBuildAuthorizationURL_MisconfiguredIDP(t *testing.T) {
	t.Parallel()

	idp := aGoogleIDP()
	idp.AuthorizationEndpoint = ""

	client := infrastructure.NewGoogleIdentityProviderClient("https://api.example.com")
	_, err := client.BuildAuthorizationURL(idp, aSessionWithState(t, "s"))
	assert.Error(t, err)
}
