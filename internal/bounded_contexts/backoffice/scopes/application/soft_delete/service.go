package soft_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
)

type Service struct {
	scopeFinder *domain_services.ScopeFinder
	repository  scopes.Repository
	authorizer  *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(scopeFinder *domain_services.ScopeFinder, repository scopes.Repository, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{scopeFinder: scopeFinder, repository: repository, authorizer: authorizer}
}

func (s *Service) SoftDelete(ctx context.Context, cmd *Command) error {
	scopeId, err := ids.NewScopeIdFromString(cmd.scopeId)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid scope id").WithCause(err, false)
	}

	scope, err := s.scopeFinder.Find(ctx, scopeId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.authorizer.Check(scope, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if scope.IsSystemResource() {
		return app_errors.ForbiddenError(scopes.ErrSystemResource).
			WithDescription("scope is a system resource and cannot be deleted")
	}

	if err = s.repository.SoftDelete(ctx, scopeId); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to soft delete scope").WithCause(err, false)
	}
	return nil
}
