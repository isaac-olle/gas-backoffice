package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
)

type ProtectedResourceFinder struct {
	repository protected_resources.Repository
}

func NewProtectedResourceFinder(repository protected_resources.Repository) *ProtectedResourceFinder {
	return &ProtectedResourceFinder{repository: repository}
}

func (f *ProtectedResourceFinder) Find(ctx context.Context, id ids.ProtectedResourceId) (*protected_resources.ProtectedResource, error) {
	pr, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("error searching protected resource").
			WithCause(err)
	}
	if pr == nil {
		return nil, app_errors.NotFoundError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("protected resource not found")
	}
	return pr, nil
}
