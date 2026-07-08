package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

type Service struct {
	repository users.Repository
	eventBus   event.Bus
}

func NewService(repository users.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	user, err := users.NewUser(cmd.id, cmd.name, cmd.ownership)
	if err != nil {
		return app_errors.BadRequestError(users.ErrUserNotCreated).
			WithDescription("failed to create user").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, user); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).
			WithDescription("failed to save user").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
