package soft_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

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

func (s *Service) SoftDelete(ctx context.Context, cmd *Command) error {
	user, err := s.finder.Find(ctx, cmd.userId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	user.Delete()

	if err = s.repository.SoftDelete(ctx, cmd.userId); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).
			WithDescription("failed to soft delete user").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, user.PullDomainEvents())

	return nil
}
