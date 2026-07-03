package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
)

type Service struct {
	repository   client_secrets.Repository
	secretHasher client_secrets.SecretHasher
	eventBus     event.Bus
}

func NewService(repository client_secrets.Repository, secretHasher client_secrets.SecretHasher, eventBus event.Bus) *Service {
	return &Service{repository: repository, secretHasher: secretHasher, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	clientId, err := ids.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("invalid client id").
			WithCause(err)
	}

	hash, err := s.secretHasher.Hash(cmd.plainSecret)
	if err != nil {
		return app_errors.InternalError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("failed to hash client secret").
			WithCause(err)
	}

	cs, err := client_secrets.NewClientSecret(cmd.id, clientId, hash)
	if err != nil {
		return app_errors.BadRequestError(client_secrets.ErrClientSecretNotCreated).
			WithDescription("failed to create client secret").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, cs); err != nil {
		return app_errors.InternalError(client_secrets.ErrClientSecretNotSaved).
			WithDescription("failed to save client secret").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, cs.PullDomainEvents())
	return nil
}
