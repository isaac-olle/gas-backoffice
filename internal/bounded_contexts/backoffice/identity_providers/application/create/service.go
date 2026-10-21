package create

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Service struct {
	repository identity_providers.Repository
	eventBus   event.Bus
	encrypter  persistance.Encrypter
}

func NewService(repository identity_providers.Repository, eventBus event.Bus, encrypter persistance.Encrypter) *Service {
	return &Service{repository: repository, eventBus: eventBus, encrypter: encrypter}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	encryptedSecret, err := s.encrypter.Encrypt(cmd.clientSecret)
	if err != nil {
		return app_errors.InternalError(identity_providers.ErrIdentityProviderNotCreated).
			WithDescription("failed to encrypt client secret").
			WithCause(fmt.Errorf("encrypt: %w", err), false)
	}

	idp, err := identity_providers.NewIdentityProvider(
		cmd.id,
		cmd.referenceName, cmd.name, cmd.issuer,
		cmd.clientID, encryptedSecret,
		cmd.authorizationEndpoint, cmd.tokenEndpoint, cmd.userinfoEndpoint, cmd.jwksURI,
		cmd.scopes,
		cmd.backofficeRedirectURI,
	)
	if err != nil {
		return app_errors.BadRequestError(identity_providers.ErrIdentityProviderNotCreated).
			WithDescription("failed to build identity provider").
			WithCause(err, true)
	}

	if err = s.repository.Save(ctx, idp); err != nil {
		return app_errors.InternalError(identity_providers.ErrIdentityProviderNotSaved).
			WithDescription("failed to save identity provider").
			WithCause(err, false)
	}

	s.eventBus.Publish(ctx, idp.PullDomainEvents())
	return nil
}
