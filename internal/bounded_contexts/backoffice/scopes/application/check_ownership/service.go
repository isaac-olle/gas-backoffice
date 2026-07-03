package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
)

type Service struct {
	repository scopes.Repository
}

func NewService(repository scopes.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CheckOwnership(ctx context.Context, qry *Query) (any, error) {
	scopeId, err := ids.NewScopeIdFromString(qry.scopeId)
	if err != nil {
		return nil, app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid scope id").
			WithCause(err)
	}

	ownership, err := scopes.NewOwnership(qry.protectedResourceId, qry.tenantId, qry.adminId)
	if err != nil {
		return nil, err
	}

	scope, err := s.repository.SearchById(ctx, scopeId)
	if err != nil {
		return nil, app_errors.NotFoundError(scopes.ErrScopeNotFound).
			WithDescription("scope not found").
			WithCause(err)
	}

	if !scope.CheckOwnership(ownership) {
		return nil, app_errors.ForbiddenError(scopes.ErrInvalidOwnership).
			WithDescription("scope does not belong to the provided protected resource")
	}

	return ownership, nil
}
