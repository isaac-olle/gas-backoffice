package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
)

type Service struct {
	repository scopes.Repository
	eventBus   event.Bus
}

func NewService(repository scopes.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	scope, err := scopes.NewScope(cmd.id, cmd.name, cmd.ownership)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotCreated).
			WithDescription("failed to create scope").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, scope); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to save scope").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, scope.PullDomainEvents())
	return nil
}
