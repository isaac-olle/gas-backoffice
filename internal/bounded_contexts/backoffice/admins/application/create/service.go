package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	repository admins.Repository
	eventBus   event.Bus
}

func NewService(repository admins.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

// Create always creates an admin as a system resource: an admin is a tenancy
// root, created out-of-band (seed, support tooling), not by another caller
// acting on it the way a tenant is created by its owning admin.
func (s *Service) Create(ctx context.Context, cmd *Command) error {
	admin, err := admins.NewAdmin(cmd.id, cmd.name, cmd.email, true)
	if err != nil {
		return app_errors.BadRequestError(admins.ErrAdminNotCreated).
			WithDescription("failed to create admin").WithCause(err, false)
	}

	if err = s.repository.Save(ctx, admin); err != nil {
		return app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("failed to save admin").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, admin.PullDomainEvents())
	return nil
}
