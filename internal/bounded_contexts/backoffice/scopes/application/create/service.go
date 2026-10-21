package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	duplicateNameChecker *scopeDomainServices.DuplicateNameChecker
	repository           scopes.Repository
	eventBus             event.Bus
}

func NewService(
	duplicateNameChecker *scopeDomainServices.DuplicateNameChecker,
	repository scopes.Repository,
	eventBus event.Bus,
) *Service {
	return &Service{duplicateNameChecker: duplicateNameChecker, repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	scope, err := scopes.NewScope(cmd.id, cmd.name, cmd.ownership, cmd.isSystem)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotCreated).
			WithDescription("failed to create scope").WithCause(err, false)
	}

	if err = s.duplicateNameChecker.Check(ctx, cmd.name, cmd.ownership.ProtectedResourceId.String()); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, scope); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to save scope").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, scope.PullDomainEvents())
	return nil
}
