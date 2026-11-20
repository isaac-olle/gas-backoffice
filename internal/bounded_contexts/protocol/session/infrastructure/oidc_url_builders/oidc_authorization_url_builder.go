package oidc_url_builders

import (
	"net/url"
	"strings"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
)

type OIDCAuthorizationURLBuilder interface {
	BuildURL(endpoint url.URL, idp *protocolShared.IdentityProviderConfig, state string) string
}

type genericOIDCAuthorizationURLBuilder struct{}

func (genericOIDCAuthorizationURLBuilder) BuildURL(endpoint url.URL, idp *protocolShared.IdentityProviderConfig, state string) string {
	endpoint.RawQuery = oidcAuthorizationQuery(endpoint, idp, state, idp.Scopes).Encode()
	return endpoint.String()
}

func oidcAuthorizationQuery(endpoint url.URL, idp *protocolShared.IdentityProviderConfig, state string, scopes []string) url.Values {
	query := endpoint.Query()
	query.Set("response_type", "code")
	query.Set("client_id", idp.ClientID)
	query.Set("redirect_uri", idp.BackofficeRedirectUri)
	query.Set("scope", strings.Join(scopes, " "))
	query.Set("state", state)
	return query
}
