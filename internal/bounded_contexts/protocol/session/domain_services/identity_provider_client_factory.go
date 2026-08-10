package domain_services

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// IdentityProviderClientFactory resolves the IdentityProviderClient adapter for
// a given provider. Like the other domain services, it holds only domain ports
// (the map value is the IdentityProviderClient interface) injected at wiring
// time, so it never imports infrastructure.
type IdentityProviderClientFactory struct {
	clients map[string]session.IdentityProviderClient
}

func NewIdentityProviderClientFactory(clients map[string]session.IdentityProviderClient) *IdentityProviderClientFactory {
	return &IdentityProviderClientFactory{clients: clients}
}

func (f *IdentityProviderClientFactory) ClientFor(idp *session.IdentityProvider) (session.IdentityProviderClient, error) {
	c, ok := f.clients[idp.Name]
	if !ok {
		return nil, app_errors.BadRequestError(session.ErrCannotAuthenticateWithIdp).
			WithDescription("unsupported identity provider: " + idp.Name)
	}
	return c, nil
}
