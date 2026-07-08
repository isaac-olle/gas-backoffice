package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
)

type Service struct {
	finder     *domain_services.UserFinder
	repository users.Repository
	eventBus   event.Bus
}

func NewService(finder *domain_services.UserFinder, repository users.Repository, eventBus event.Bus) *Service {
	return &Service{finder: finder, repository: repository, eventBus: eventBus}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	userId, err := ids.NewUserIdFromString(cmd.userId)
	if err != nil {
		return app_errors.BadRequestError(users.ErrUserNotFound).
			WithDescription("invalid user id").
			WithCause(err)
	}

	user, err := s.finder.Find(ctx, userId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	user.Delete()

	if err = s.repository.HardDelete(ctx, userId); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).
			WithDescription("failed to hard delete user").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
