package add_allowed_client

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

type Service struct {
	scopeFinder            *domain_services.ScopeFinder
	repository             scopes.Repository
	eventBus               event.Bus
	clientProjectionFinder *domain_services.ClientProjectionFinder
	authorizer             *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(scopeFinder *domain_services.ScopeFinder, repository scopes.Repository, eventBus event.Bus, clientProjectionFinder *domain_services.ClientProjectionFinder, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{scopeFinder: scopeFinder, repository: repository, eventBus: eventBus, clientProjectionFinder: clientProjectionFinder, authorizer: authorizer}
}

func (s *Service) AddAllowedClient(ctx context.Context, cmd *Command) error {
	scopeId, err := ids2.NewScopeIdFromString(cmd.scopeId)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid scope id").WithCause(err, false)
	}

	clientId, err := ids2.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid client id").WithCause(err, false)
	}

	scope, err := s.scopeFinder.Find(ctx, scopeId)
	if err != nil {
		return app_errors.AddTrace(err)
	}
	if err = s.authorizer.Check(scope, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	client, err := s.clientProjectionFinder.Find(ctx, clientId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = scope.EnsureClientSubscription(client.SubscriptionId); err != nil {
		return app_errors.AddTrace(err)
	}

	scope.AddAllowedClient(clientId)

	if err = s.repository.Save(ctx, scope); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(scopes.ErrScopeConcurrentModification).
				WithDescription("scope was modified by another request; reload and retry").
				WithCause(err, false)
		}
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to save scope").WithCause(

			err, false)

	}

	s.eventBus.Publish(ctx, scope.PullDomainEvents())
	return nil
}
