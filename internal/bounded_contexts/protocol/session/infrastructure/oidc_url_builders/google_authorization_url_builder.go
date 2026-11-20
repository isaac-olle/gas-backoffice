package oidc_url_builders

import (
	"net/url"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
)

type GoogleAuthorizationURLBuilder struct{}

func (GoogleAuthorizationURLBuilder) BuildURL(endpoint url.URL, idp *protocolShared.IdentityProviderConfig, state string) string {
	query := oidcAuthorizationQuery(endpoint, idp, state, idp.Scopes)
	query.Set("prompt", "consent select_account")
	query.Set("access_type", "offline")
	endpoint.RawQuery = query.Encode()
	return endpoint.String()
}
