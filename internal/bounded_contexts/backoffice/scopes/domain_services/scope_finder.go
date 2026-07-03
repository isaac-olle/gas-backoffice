package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
)

type ScopeFinder struct {
	repository scopes.Repository
}

func NewScopeFinder(repository scopes.Repository) *ScopeFinder {
	return &ScopeFinder{repository: repository}
}

func (f *ScopeFinder) Find(ctx context.Context, id ids.ScopeId) (*scopes.Scope, error) {
	return f.repository.SearchById(ctx, id)
}
