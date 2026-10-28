package add_redirect_uri

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

type Service struct {
	clientFinder *domain_services.ClientFinder
	authorizer   *sharedDomainServices.SubscriptionAuthorizer
	repository   clients.Repository
	eventBus     event.Bus
}

func NewService(clientFinder *domain_services.ClientFinder, repository clients.Repository, eventBus event.Bus, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{clientFinder: clientFinder, repository: repository, eventBus: eventBus, authorizer: authorizer}
}

func (s *Service) Add(ctx context.Context, cmd *Command) error {
	uri, err := vos.NewRedirectUri(cmd.redirectUri)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	client, err := s.clientFinder.Find(ctx, cmd.clientId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.authorizer.Check(client, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if client.IsSystemResource() {
		return app_errors.ForbiddenError(clients.ErrSystemResource).
			WithDescription("client is a system resource and cannot be modified")
	}

	client.AddRedirectUri(uri)

	if err = s.repository.Save(ctx, client); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(clients.ErrClientConcurrentModification).
				WithDescription("client was modified by another request; reload and retry").
				WithCause(err, false)
		}
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
