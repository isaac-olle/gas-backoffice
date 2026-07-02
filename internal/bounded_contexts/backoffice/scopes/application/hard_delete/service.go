package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/domain_services"
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
		return err
	}

	scope.Delete()

	if err := s.repository.HardDelete(ctx, cmd.scopeId); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithCause(err)
	}

	s.eventBus.Publish(ctx, scope.PullDomainEvents())
	return nil
}
