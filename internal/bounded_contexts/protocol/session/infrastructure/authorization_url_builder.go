package infrastructure

import (
	"fmt"
	"net/url"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/oidc_url_builders"
)

// AuthorizationURLBuilder builds the front-channel authorization URL for an
// OIDC provider following RFC 6749 / OpenID Connect Core.
type AuthorizationURLBuilder struct {
	factory *oidc_url_builders.OIDCAuthorizationURLBuilderFactory
}

func NewAuthorizationURLBuilder(factory *oidc_url_builders.OIDCAuthorizationURLBuilderFactory) *AuthorizationURLBuilder {
	return &AuthorizationURLBuilder{factory: factory}
}

func (b *AuthorizationURLBuilder) BuildAuthorizationURL(idp *protocolShared.IdentityProviderConfig, state string) (string, error) {
	if idp.AuthorizationEndpoint == "" || idp.ClientID == "" {
		return "", fmt.Errorf("identity provider %q is missing authorization endpoint or client id", idp.Name)
	}

	endpoint, err := url.Parse(idp.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid authorization endpoint for identity provider %q: %w", idp.Name, err)
	}

	builder := b.factory.ForHost(endpoint.Hostname())
	return builder.BuildURL(*endpoint, idp, state), nil
}
