package authenticate_with_idp

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Service starts the federated authentication: it takes the (already created)
// session and the identity provider, resolves the right IdentityProviderClient
// adapter via the factory, and builds the authorization URL the browser is
// redirected to. Only Google is registered in the factory for now.
type Service struct {
	sessionFinder *domain_services.SessionFinder
	idpFinder     *domain_services.IdentityProviderFinder
	clientFactory *domain_services.IdentityProviderClientFactory
}

func NewService(
	sessionFinder *domain_services.SessionFinder,
	idpFinder *domain_services.IdentityProviderFinder,
	clientFactory *domain_services.IdentityProviderClientFactory,
) *Service {
	return &Service{sessionFinder: sessionFinder, idpFinder: idpFinder, clientFactory: clientFactory}
}

func (s *Service) Authenticate(ctx context.Context, qry *Query) (any, error) {
	sess, err := s.sessionFinder.FindById(ctx, qry.sessionID)
	if err != nil {
		return nil, err
	}

	idp, err := s.idpFinder.Find(ctx, qry.idpID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientFactory.ClientFor(idp)
	if err != nil {
		return nil, err
	}

	authorizationURL, err := client.BuildAuthorizationURL(idp, sess)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrCannotAuthenticateWithIdp).
			WithDescription("failed to build identity provider authorization request").
			WithCause(err)
	}

	return authorizationURL, nil
}
