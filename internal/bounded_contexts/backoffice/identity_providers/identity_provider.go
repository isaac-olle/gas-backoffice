package identity_providers

import (
	"net/url"
	"strings"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// IdentityProvider is a read-only projection of a configured external OIDC
// provider (row in backoffice.oauth_identity_providers). It carries no domain
// events: identity providers are seeded/configured in the DB, never mutated
// through the domain. Fields stay private; behaviour is exposed via methods so
// the aggregate never leaks its internals (no getters).
type IdentityProvider struct {
	id                    ids.IdentityProviderId
	name                  string
	issuer                string
	clientID              string
	clientSecret          string
	authorizationEndpoint string
	tokenEndpoint         string
	userinfoEndpoint      string
	jwksURI               string
	scopes                []string
	enabled               bool
}

func FromPrimitives(p IdentityProviderPrimitives) (*IdentityProvider, error) {
	id, err := ids.NewIdentityProviderIdFromString(p.Id)
	if err != nil {
		return nil, err
	}
	return &IdentityProvider{
		id:                    id,
		name:                  p.Name,
		issuer:                p.Issuer,
		clientID:              p.ClientId,
		clientSecret:          p.ClientSecret,
		authorizationEndpoint: p.AuthorizationEndpoint,
		tokenEndpoint:         p.TokenEndpoint,
		userinfoEndpoint:      p.UserinfoEndpoint,
		jwksURI:               p.JwksUri,
		scopes:                p.Scopes,
		enabled:               p.Enabled,
	}, nil
}

func (p *IdentityProvider) IsEnabled() bool { return p.enabled }

// AuthorizationURL builds the front-channel URL the browser is redirected to.
// It never exposes the client secret (that is only used at the token exchange).
func (p *IdentityProvider) AuthorizationURL(redirectURI, state, nonce string) (string, error) {
	if p.authorizationEndpoint == "" || p.clientID == "" {
		return "", app_errors.InternalError(ErrIdentityProviderMisconfigured).
			WithDescription("identity provider is missing authorization endpoint or client id")
	}
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", p.clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", strings.Join(p.scopes, " "))
	q.Set("state", state)
	q.Set("nonce", nonce)
	return p.authorizationEndpoint + "?" + q.Encode(), nil
}
