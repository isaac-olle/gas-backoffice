package oidc_url_builders

import (
	"net/url"
	"slices"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
)

type MicrosoftAuthorizationURLBuilder struct{}

func (MicrosoftAuthorizationURLBuilder) BuildURL(endpoint url.URL, idp *protocolShared.IdentityProviderConfig, state string) string {
	scopes := idp.Scopes
	if !slices.Contains(scopes, "offline_access") {
		scopes = append(slices.Clone(scopes), "offline_access")
	}

	query := oidcAuthorizationQuery(endpoint, idp, state, scopes)
	query.Set("prompt", "select_account")
	endpoint.RawQuery = query.Encode()
	return endpoint.String()
}
