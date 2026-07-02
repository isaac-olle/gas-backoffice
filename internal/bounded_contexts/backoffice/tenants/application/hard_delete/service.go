package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

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

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	tenantId, err := ids.NewTenantIdFromString(cmd.tenantId)
	if err != nil {
		return app_errors.BadRequestError(tenants.ErrTenantNotFound).
			WithDescription("invalid tenant id").
			WithCause(err)
	}

	tenant, err := s.finder.Find(ctx, tenantId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	tenant.Delete()

	if err = s.repository.HardDelete(ctx, tenantId); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).
			WithDescription("failed to hard delete tenant").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
