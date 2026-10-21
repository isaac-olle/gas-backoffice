package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ScopeFinder struct {
	repository scopes.Repository
}

func NewScopeFinder(repository scopes.Repository) *ScopeFinder {
	return &ScopeFinder{repository: repository}
}

func (f *ScopeFinder) Find(ctx context.Context, id ids.ScopeId) (*scopes.Scope, error) {
	scope, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(scopes.ErrScopeNotFound).
			WithDescription("error searching scope").WithCause(err, false)
	}
	if scope == nil {
		return nil, app_errors.NotFoundError(scopes.ErrScopeNotFound).
			WithDescription("scope not found")
	}
	return scope, nil
}
