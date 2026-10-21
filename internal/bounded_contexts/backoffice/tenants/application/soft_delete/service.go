package soft_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
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
	tenant, err := s.finder.Find(ctx, cmd.userId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = tenant.EnsureTenantAccess(cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if tenant.IsSystemResource() {
		return app_errors.ForbiddenError(tenants.ErrSystemResource).
			WithDescription("tenant is a system resource and cannot be deleted")
	}

	tenant.Delete()

	if err = s.repository.SoftDelete(ctx, cmd.userId); err != nil {
		return app_errors.InternalError(tenants.ErrUserNotSaved).
			WithDescription("failed to soft delete tenant").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())

	return nil
}
