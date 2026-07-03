package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
)

type Service struct {
	repository tenants.Repository
	eventBus   event.Bus
}

func NewService(repository tenants.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	tenant, err := tenants.NewTenant(cmd.id, cmd.name, cmd.ownership)
	if err != nil {
		return app_errors.BadRequestError(tenants.ErrTenantNotCreated).
			WithDescription("failed to create tenant").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, tenant); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).
			WithDescription("failed to save tenant").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
