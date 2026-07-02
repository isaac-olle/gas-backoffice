package soft_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/domain_services"
)

type Service struct {
	finder     *domain_services.TenantFinder
	repository tenants.Repository
	eventBus   event.Bus
}

func NewService(finder *domain_services.TenantFinder, repository tenants.Repository, eventBus event.Bus) *Service {
	return &Service{finder: finder, repository: repository, eventBus: eventBus}
}

func (s *Service) SoftDelete(ctx context.Context, cmd *Command) error {
	tenant, err := s.finder.Find(ctx, cmd.tenantId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	tenant.Delete()

	if err = s.repository.SoftDelete(ctx, cmd.tenantId); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).
			WithDescription("failed to soft delete tenant").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())

	return nil
}
