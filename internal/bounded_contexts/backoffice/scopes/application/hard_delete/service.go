package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
)

type Service struct {
	scopeFinder *domain_services.ScopeFinder
	repository  scopes.Repository
	eventBus    event.Bus
}

func NewService(scopeFinder *domain_services.ScopeFinder, repository scopes.Repository, eventBus event.Bus) *Service {
	return &Service{scopeFinder: scopeFinder, repository: repository, eventBus: eventBus}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	scope, err := s.scopeFinder.Find(ctx, cmd.scopeId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	subscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if err = scope.EnsureOwnership(cmd.protectedResourceId, []ids.SubscriptionId{subscriptionId}); err != nil {
		return app_errors.AddTrace(err)
	}

	if scope.IsSystemResource() {
		return app_errors.ForbiddenError(scopes.ErrSystemResource).
			WithDescription("scope is a system resource and cannot be deleted")
	}

	scope.Delete()

	if err := s.repository.HardDelete(ctx, cmd.scopeId); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).WithCause(
			err, false)

	}

	s.eventBus.Publish(ctx, scope.PullDomainEvents())
	return nil
}
