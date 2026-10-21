package authenticate_client

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type PublicClientMethod struct {
	clientFinder *domain_services.ClientFinder
}

func NewPublicClientMethod(clientFinder *domain_services.ClientFinder) *PublicClientMethod {
	return &PublicClientMethod{clientFinder: clientFinder}
}

func (m *PublicClientMethod) CanProcess(creds Credentials) bool {
	return creds.IsFormEncoded && creds.ClientID != "" && creds.ClientSecret == ""
}

func (m *PublicClientMethod) Authenticate(ctx context.Context, creds Credentials) (Result, error) {
	client, err := m.clientFinder.Find(ctx, creds.ClientID)
	if err != nil {
		return Result{}, err
	}
	if !client.Public {
		return Result{}, app_errors.ForbiddenError(session.ErrClientNotPublic).
			WithDescription("client is not a public client")
	}
	return Result{ClientID: creds.ClientID, IsPublic: true}, nil
}
