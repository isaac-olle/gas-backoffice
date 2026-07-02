package change_name

import (
	"context"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/vos"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/domain_services"
)

type Service struct {
	tenantFinder *domain_services.TenantFinder
	repository   tenants.Repository
	eventBus     event.Bus
}

func NewService(tenantFinder *domain_services.TenantFinder, repository tenants.Repository, eventBus event.Bus) *Service {
	return &Service{tenantFinder: tenantFinder, repository: repository, eventBus: eventBus}
}

func (s *Service) ChangeName(ctx context.Context, cmd *Command) error {
	newName, err := vos.NewTenantName(cmd.name)
	if err != nil {
		return app_errors.BadRequestError(tenants.ErrTenantNameNotChanged).
			WithDescription("invalid tenant name").
			WithCause(err)
	}

	tenant, err := s.tenantFinder.Find(ctx, cmd.tenantId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = tenant.ChangeName(newName); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, tenant); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).
			WithDescription("failed to save tenant").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
