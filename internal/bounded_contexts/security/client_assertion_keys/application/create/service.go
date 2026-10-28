package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
)

type Service struct {
	clientProjectionRepository client_projection.Repository
	repository                 client_assertion_keys.Repository
	eventBus                   event.Bus
	authorizer                 *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(
	clientProjectionRepository client_projection.Repository,
	repository client_assertion_keys.Repository,
	eventBus event.Bus,
	authorizer *sharedDomainServices.SubscriptionAuthorizer,
) *Service {
	return &Service{clientProjectionRepository: clientProjectionRepository, repository: repository, eventBus: eventBus, authorizer: authorizer}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	projection, err := s.clientProjectionRepository.Find(ctx, cmd.clientId)
	if err != nil {
		return app_errors.NotFoundError(client_assertion_keys.ErrClientAssertionKeyNotCreated).
			WithDescription("client not found").WithCause(err, false)
	}

	if err = s.authorizer.Check(projection, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	// The new key is owned by the CLIENT's owner (not necessarily the caller,
	// who may be the owning admin acting on a tenant's client).
	subscriptionId, err := ids.NewSubscriptionIdFromString(projection.SubscriptionId)
	if err != nil {
		return app_errors.InternalError(client_assertion_keys.ErrClientAssertionKeyNotCreated).
			WithDescription("invalid client subscription id").WithCause(err, false)
	}

	cak, err := client_assertion_keys.NewClientAssertionKeyForSubscription(cmd.id, cmd.clientId, subscriptionId, cmd.algorithm, cmd.publicKey)
	if err != nil {
		return app_errors.BadRequestError(client_assertion_keys.ErrClientAssertionKeyNotCreated).
			WithDescription("failed to create client assertion key").WithCause(

			err, false)

	}

	if err = s.repository.Save(ctx, cak); err != nil {
		return app_errors.InternalError(client_assertion_keys.ErrClientAssertionKeyNotSaved).
			WithDescription("failed to save client assertion key").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, cak.PullDomainEvents())
	return nil
}
