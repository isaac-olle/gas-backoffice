package infrastructure

import (
	"fmt"
	"net/url"
	"strings"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
)

// AuthorizationURLBuilder builds the front-channel authorization URL for an
// OIDC provider following RFC 6749 / OpenID Connect Core.
type AuthorizationURLBuilder struct{}

func NewAuthorizationURLBuilder() *AuthorizationURLBuilder {
	return &AuthorizationURLBuilder{}
}

func (b *AuthorizationURLBuilder) BuildAuthorizationURL(idp *protocolShared.IdentityProviderConfig, state string) (string, error) {
	if idp.AuthorizationEndpoint == "" || idp.ClientID == "" {
		return "", fmt.Errorf("identity provider %q is missing authorization endpoint or client id", idp.Name)
	}

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", idp.ClientID)
	q.Set("redirect_uri", idp.BackofficeRedirectUri)
	q.Set("scope", strings.Join(idp.Scopes, " "))
	q.Set("state", state)
	// access_type=offline is required for Google to hand back a refresh_token;
	// prompt=consent forces the consent screen every time so that refresh_token
	// is actually included (Google only sends it the first time a user consents
	// otherwise), and also lets users pick another account instead of silently
	// reusing the browser's already-signed-in one.
	q.Set("access_type", "offline")
	q.Set("prompt", "consent select_account")

	return idp.AuthorizationEndpoint + "?" + q.Encode(), nil
}
