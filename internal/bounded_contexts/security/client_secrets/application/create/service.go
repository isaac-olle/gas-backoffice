package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Service struct {
	clientProjectionRepository client_projection.Repository
	repository                 client_secrets.Repository
	secretHasher               client_secrets.SecretHasher
	eventBus                   event.Bus
}

func NewService(
	clientProjectionRepository client_projection.Repository,
	repository client_secrets.Repository,
	secretHasher client_secrets.SecretHasher,
	eventBus event.Bus,
) *Service {
	return &Service{clientProjectionRepository: clientProjectionRepository, repository: repository, secretHasher: secretHasher, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	projection, err := s.clientProjectionRepository.Find(ctx, cmd.clientId)
	if err != nil {
		return app_errors.NotFoundError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("client not found").WithCause(err, false)
	}

	pathSubscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if !projection.IsAccessibleBy([]string{pathSubscriptionId.String()}) {
		return app_errors.ForbiddenError(client_secrets.ErrInvalidOwnership).
			WithDescription("client does not belong to the provided subscription")
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
