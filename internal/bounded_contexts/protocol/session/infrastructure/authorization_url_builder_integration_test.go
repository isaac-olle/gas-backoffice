//go:build integration

package infrastructure_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/oidc_url_builders"
)

func aGoogleIDP() *protocolShared.IdentityProviderConfig {
	return protocolShared.NewIdentityProviderConfig(
		"idp-id", "Google", "https://accounts.google.com", "client-abc", "",
		"https://accounts.google.com/o/oauth2/v2/auth", "", "https://www.googleapis.com/oauth2/v3/certs",
		[]string{"openid", "email"},
		"https://api.example.com/v1/auth/callback/google",
	)
}

func TestIntegrationBuildAuthorizationURL(t *testing.T) {
	t.Parallel()

	client := infrastructure.NewAuthorizationURLBuilder(oidc_url_builders.NewOIDCAuthorizationURLBuilderFactory())
	raw, err := client.BuildAuthorizationURL(aGoogleIDP(), "the-state")
	require.NoError(t, err)

	u, err := url.Parse(raw)
	require.NoError(t, err)
	assert.Equal(t, "accounts.google.com", u.Host)

	q := u.Query()
	assert.Equal(t, "code", q.Get("response_type"))
	assert.Equal(t, "client-abc", q.Get("client_id"))
	assert.Equal(t, "https://api.example.com/v1/auth/callback/google", q.Get("redirect_uri"))
	assert.Equal(t, "openid email", q.Get("scope"))
	assert.Equal(t, "the-state", q.Get("state"))
	assert.Equal(t, "offline", q.Get("access_type"))
	assert.Equal(t, "consent select_account", q.Get("prompt"))
}

func TestIntegrationBuildAuthorizationURL_MisconfiguredIDP(t *testing.T) {
	t.Parallel()

	idp := aGoogleIDP()
	idp.AuthorizationEndpoint = ""

	client := infrastructure.NewAuthorizationURLBuilder(oidc_url_builders.NewOIDCAuthorizationURLBuilderFactory())
	_, err := client.BuildAuthorizationURL(idp, "s")
	assert.Error(t, err)
}

func TestIntegrationBuildAuthorizationURL_ProviderParameters(t *testing.T) {
	t.Parallel()

	client := infrastructure.NewAuthorizationURLBuilder(oidc_url_builders.NewOIDCAuthorizationURLBuilderFactory())
	for _, tc := range []struct {
		name     string
		endpoint string
		scopes   []string
		scope    string
		prompt   string
	}{
		{"Microsoft", "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize", []string{"openid", "email"}, "openid email offline_access", "select_account"},
		{"Microsoft with offline access", "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize", []string{"openid", "offline_access"}, "openid offline_access", "select_account"},
		{"Generic OIDC", "https://idp.example.com/authorize", []string{"openid", "email"}, "openid email", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			idp := aGoogleIDP()
			idp.AuthorizationEndpoint = tc.endpoint
			idp.Scopes = tc.scopes

			raw, err := client.BuildAuthorizationURL(idp, "the-state")
			require.NoError(t, err)
			u, err := url.Parse(raw)
			require.NoError(t, err)
			q := u.Query()
			assert.Equal(t, tc.scope, q.Get("scope"))
			assert.Equal(t, tc.prompt, q.Get("prompt"))
			assert.NotContains(t, q, "access_type")
		})
	}
}
