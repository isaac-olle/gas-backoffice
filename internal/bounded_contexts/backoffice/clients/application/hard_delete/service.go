package hard_delete

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
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

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	clientId, err := ids.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id").WithCause(err, false)
	}

	client, err := s.clientFinder.Find(ctx, clientId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.authorizer.Check(client, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if client.IsSystemResource() {
		return app_errors.ForbiddenError(clients.ErrSystemResource).
			WithDescription("client is a system resource and cannot be deleted")
	}

	client.Delete()

	if err = s.repository.HardDelete(ctx, clientId); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to hard delete client").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
