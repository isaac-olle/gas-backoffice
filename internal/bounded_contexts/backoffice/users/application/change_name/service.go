package change_name

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
)

type Service struct {
	userFinder *domain_services.UserFinder
	repository users.Repository
	eventBus   event.Bus
}

func NewService(userFinder *domain_services.UserFinder, repository users.Repository, eventBus event.Bus) *Service {
	return &Service{userFinder: userFinder, repository: repository, eventBus: eventBus}
}

func (s *Service) ChangeName(ctx context.Context, cmd *Command) error {
	newName, err := vos.NewUserName(cmd.name)
	if err != nil {
		return app_errors.BadRequestError(users.ErrUserNameNotChanged).
			WithDescription("invalid user name").WithCause(

			err, false)

	}

	user, err := s.userFinder.Find(ctx, cmd.userId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = user.EnsureUserAccess(cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if user.IsSystemResource() {
		return app_errors.ForbiddenError(users.ErrSystemResource).
			WithDescription("user is a system resource and cannot be modified")
	}

	if err = user.ChangeName(newName); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, user); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(users.ErrUserConcurrentModification).
				WithDescription("user was modified by another request; reload and retry").
				WithCause(err, false)
		}
		return app_errors.InternalError(users.ErrUserNotSaved).
			WithDescription("failed to save user").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
