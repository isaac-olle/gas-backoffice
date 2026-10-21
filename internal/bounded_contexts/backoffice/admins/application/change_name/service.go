package change_name

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
)

type Service struct {
	adminFinder *domain_services.AdminFinder
	repository  admins.Repository
	eventBus    event.Bus
}

func NewService(adminFinder *domain_services.AdminFinder, repository admins.Repository, eventBus event.Bus) *Service {
	return &Service{adminFinder: adminFinder, repository: repository, eventBus: eventBus}
}

func (s *Service) ChangeName(ctx context.Context, cmd *Command) error {
	newName, err := vos.NewUserName(cmd.name)
	if err != nil {
		return app_errors.BadRequestError(admins.ErrAdminNameNotChanged).
			WithDescription("invalid admin name").WithCause(err, false)
	}

	admin, err := s.adminFinder.Find(ctx, cmd.adminId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = admin.EnsureSelf(cmd.callerId); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = admin.ChangeName(newName); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, admin); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(admins.ErrAdminConcurrentModification).
				WithDescription("admin was modified by another request; reload and retry").
				WithCause(err, false)
		}
		return app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("failed to save admin").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, admin.PullDomainEvents())
	return nil
}
