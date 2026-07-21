package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys"
)

type Service struct {
	repository user_assertion_keys.Repository
	eventBus   event.Bus
}

func NewService(repository user_assertion_keys.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	userId, err := ids.NewUserIdFromString(cmd.userId)
	if err != nil {
		return app_errors.BadRequestError(user_assertion_keys.ErrUserAssertionKeyNotCreated).
			WithDescription("invalid user id").
			WithCause(err)
	}

	uak, err := user_assertion_keys.NewUserAssertionKey(cmd.id, userId, cmd.algorithm, cmd.publicKey)
	if err != nil {
		return app_errors.BadRequestError(user_assertion_keys.ErrUserAssertionKeyNotCreated).
			WithDescription("failed to create user assertion key").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, uak); err != nil {
		return app_errors.InternalError(user_assertion_keys.ErrUserAssertionKeyNotSaved).
			WithDescription("failed to save user assertion key").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, uak.PullDomainEvents())
	return nil
}
