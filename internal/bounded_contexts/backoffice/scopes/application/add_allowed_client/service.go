package add_allowed_client

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
)

type Service struct {
	repository scopes.Repository
	eventBus   event.Bus
}

func NewService(repository scopes.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) AddAllowedClient(ctx context.Context, cmd *Command) error {
	scopeId, err := ids2.NewScopeIdFromString(cmd.scopeId)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid scope id").
			WithCause(err)
	}

	clientId, err := ids2.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid client id").
			WithCause(err)
	}

	scope, err := s.repository.SearchById(ctx, scopeId)
	if err != nil {
		return app_errors.NotFoundError(scopes.ErrScopeNotFound).
			WithDescription("scope not found").
			WithCause(err)
	}

	scope.AddAllowedClient(clientId)

	if err = s.repository.Save(ctx, scope); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to save scope").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, scope.PullDomainEvents())
	return nil
}
