package create_tenant

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	repository tenants.Repository
	eventBus   event.Bus
}

func NewService(repository tenants.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	tenant, err := tenants.NewTenant(cmd.id, cmd.name, cmd.adminId, "", cmd.isSystem)
	if err != nil {
		return app_errors.BadRequestError(tenants.ErrUserNotCreated).
			WithDescription("failed to create tenant").WithCause(err, false)
	}

	if err = s.repository.Save(ctx, tenant); err != nil {
		return app_errors.InternalError(tenants.ErrUserNotSaved).
			WithDescription("failed to save tenant").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
