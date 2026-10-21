package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

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

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	userId, err := ids.NewUserIdFromString(cmd.userId)
	if err != nil {
		return app_errors.BadRequestError(tenants.ErrUserNotFound).
			WithDescription("invalid tenant id").WithCause(

			err, false)

	}

	tenant, err := s.finder.Find(ctx, userId)
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

	if err = s.repository.HardDelete(ctx, userId); err != nil {
		return app_errors.InternalError(tenants.ErrUserNotSaved).
			WithDescription("failed to hard delete tenant").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
