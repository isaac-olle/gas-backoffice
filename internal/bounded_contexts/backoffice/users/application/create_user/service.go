package create_user

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

type Service struct {
	repository users.Repository
	eventBus   event.Bus
}

func NewService(repository users.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	role := vo.SubscriptionUserRole
	if cmd.isSystem {
		role = vo.SystemAdministratorRole
	}
	user, err := users.NewUser(cmd.userId, cmd.name, cmd.email, role, cmd.isSystem)
	if err != nil {
		return app_errors.BadRequestError(users.ErrUserNotCreated).
			WithDescription("failed to create user").WithCause(err, false)
	}

	if err = s.repository.Save(ctx, user); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).
			WithDescription("failed to save user").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
