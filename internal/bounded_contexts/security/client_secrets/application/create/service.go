package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

type Service struct {
	clientProjectionRepository client_projection.Repository
	repository                 client_secrets.Repository
	secretHasher               client_secrets.SecretHasher
	eventBus                   event.Bus
	authorizer                 *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(
	clientProjectionRepository client_projection.Repository,
	repository client_secrets.Repository,
	secretHasher client_secrets.SecretHasher,
	eventBus event.Bus,
	authorizer *sharedDomainServices.SubscriptionAuthorizer,
) *Service {
	return &Service{clientProjectionRepository: clientProjectionRepository, repository: repository, secretHasher: secretHasher, eventBus: eventBus, authorizer: authorizer}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	projection, err := s.clientProjectionRepository.Find(ctx, cmd.clientId)
	if err != nil {
		return app_errors.NotFoundError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("client not found").WithCause(err, false)
	}

	if err = s.authorizer.Check(projection, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	// The new secret is owned by the CLIENT's owner (not necessarily the caller,
	// who may be the owning admin acting on a tenant's client).
	subscriptionId, err := ids.NewSubscriptionIdFromString(projection.SubscriptionId)
	if err != nil {
		return app_errors.InternalError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("invalid client subscription id").WithCause(err, false)
	}

	cs, err := client_secrets.NewClientSecretForSubscription(cmd.id, cmd.clientId, subscriptionId, cmd.plainSecret, s.secretHasher)
	if err != nil {
		return app_errors.BadRequestError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("failed to create client secret").WithCause(

			err, false)

	}

	if err = s.repository.Save(ctx, cs); err != nil {
		return app_errors.InternalError(client_secrets.ErrClientSecretNotSaved).
			WithDescription("failed to save client secret").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, cs.PullDomainEvents())
	return nil
}
