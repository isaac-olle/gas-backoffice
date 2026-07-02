package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys"
)

type Service struct {
	repository client_assertion_keys.Repository
	eventBus   event.Bus
}

func NewService(repository client_assertion_keys.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	cak, err := client_assertion_keys.NewClientAssertionKey(cmd.id, cmd.clientId, cmd.algorithm, cmd.publicKey)
	if err != nil {
		return app_errors.BadRequestError(client_assertion_keys.ErrClientAssertionKeyNotCreated).
			WithDescription("failed to create client assertion key").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, cak); err != nil {
		return app_errors.InternalError(client_assertion_keys.ErrClientAssertionKeyNotSaved).
			WithDescription("failed to save client assertion key").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, cak.PullDomainEvents())
	return nil
}
