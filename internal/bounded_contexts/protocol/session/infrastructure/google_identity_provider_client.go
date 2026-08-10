package infrastructure

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

const callbackPath = "/v1/auth/callback/"

// GoogleIdentityProviderClient talks to Google as an OIDC provider. For now it
// only builds the front-channel authorization URL; the back-channel calls
// (token exchange, JWKS verification) will live here too when the callback is
// implemented. The client secret is never exposed in the authorization URL.
type GoogleIdentityProviderClient struct {
	publicBaseURL string
}

func NewGoogleIdentityProviderClient(publicBaseURL string) *GoogleIdentityProviderClient {
	return &GoogleIdentityProviderClient{publicBaseURL: publicBaseURL}
}

func (g *GoogleIdentityProviderClient) BuildAuthorizationURL(idp *session.IdentityProvider, sess *session.AuthorizationCodeSession) (string, error) {
	if idp.AuthorizationEndpoint == "" || idp.ClientID == "" {
		return "", fmt.Errorf("identity provider %q is missing authorization endpoint or client id", idp.Name)
	}

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", idp.ClientID)
	q.Set("redirect_uri", g.publicBaseURL+callbackPath+idp.ID)
	q.Set("scope", strings.Join(idp.Scopes, " "))
	q.Set("state", sess.ToPrimitives().State)

	return idp.AuthorizationEndpoint + "?" + q.Encode(), nil
}
