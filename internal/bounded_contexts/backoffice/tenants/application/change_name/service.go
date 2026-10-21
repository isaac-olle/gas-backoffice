package change_name

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
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
	newName, err := vos.NewUserName(cmd.name)
	if err != nil {
		return app_errors.BadRequestError(tenants.ErrUserNameNotChanged).
			WithDescription("invalid tenant name").WithCause(

			err, false)

	}

	tenant, err := s.tenantFinder.Find(ctx, cmd.userId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = tenant.EnsureTenantAccess(cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if tenant.IsSystemResource() {
		return app_errors.ForbiddenError(tenants.ErrSystemResource).
			WithDescription("tenant is a system resource and cannot be modified")
	}

	if err = tenant.ChangeName(newName); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, tenant); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(tenants.ErrUserConcurrentModification).
				WithDescription("tenant was modified by another request; reload and retry").
				WithCause(err, false)
		}
		return app_errors.InternalError(tenants.ErrUserNotSaved).
			WithDescription("failed to save tenant").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
