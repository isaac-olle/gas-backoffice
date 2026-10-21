package identity_providers_test

import (
	"net/url"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validPrimitives(id string) identity_providers.IdentityProviderPrimitives {
	return identity_providers.IdentityProviderPrimitives{
		Id:                    id,
		ReferenceName:         "google",
		Name:                  "google",
		Issuer:                "https://accounts.google.com",
		ClientId:              "client-123",
		ClientSecret:          "secret-xyz",
		AuthorizationEndpoint: "https://accounts.google.com/o/oauth2/v2/auth",
		TokenEndpoint:         "https://oauth2.googleapis.com/token",
		JwksUri:               "https://www.googleapis.com/oauth2/v3/certs",
		Scopes:                []string{"openid", "email", "profile"},
		BackofficeRedirectUri: "https://api.example.com/v1/auth/callback/x",
		Enabled:               true,
	}
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs an enabled provider", func(t *testing.T) {
		t.Parallel()
		idp, err := identity_providers.FromPrimitives(validPrimitives(ids2.NewIdentityProviderId().String()))
		require.NoError(t, err)
		assert.True(t, idp.IsEnabled())
	})

	t.Run("rejects invalid id", func(t *testing.T) {
		t.Parallel()
		_, err := identity_providers.FromPrimitives(validPrimitives("bad-id"))
		assert.Error(t, err)
	})
}

func TestAuthorizationURL(t *testing.T) {
	t.Parallel()

	t.Run("builds the authorization url with the expected query params", func(t *testing.T) {
		t.Parallel()
		idp, err := identity_providers.FromPrimitives(validPrimitives(ids2.NewIdentityProviderId().String()))
		require.NoError(t, err)

		raw, err := idp.AuthorizationURL("https://api.example.com/v1/auth/callback/x", "the-state", "the-nonce")
		require.NoError(t, err)

		u, err := url.Parse(raw)
		require.NoError(t, err)
		assert.Equal(t, "accounts.google.com", u.Host)
		q := u.Query()
		assert.Equal(t, "code", q.Get("response_type"))
		assert.Equal(t, "client-123", q.Get("client_id"))
		assert.Equal(t, "https://api.example.com/v1/auth/callback/x", q.Get("redirect_uri"))
		assert.Equal(t, "openid email profile", q.Get("scope"))
		assert.Equal(t, "the-state", q.Get("state"))
		assert.Equal(t, "the-nonce", q.Get("nonce"))
		assert.NotContains(t, raw, "secret-xyz")
	})

	t.Run("fails when misconfigured (missing authorization endpoint)", func(t *testing.T) {
		t.Parallel()
		p := validPrimitives(ids2.NewIdentityProviderId().String())
		p.AuthorizationEndpoint = ""
		_, err := identity_providers.FromPrimitives(p)
		assert.Error(t, err)
	})
}
