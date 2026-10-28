package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ClientProjectionFinder struct {
	repository scopes.ClientProjectionRepository
}

func NewClientProjectionFinder(repository scopes.ClientProjectionRepository) *ClientProjectionFinder {
	return &ClientProjectionFinder{repository: repository}
}

func (f *ClientProjectionFinder) Find(ctx context.Context, id ids.ClientId) (*scopes.ClientProjection, error) {
	projection, err := f.repository.FindClientProjection(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError("CLIENT_PROJECTION_NOT_FOUND").WithCause(err, false)
	}
	if projection == nil {
		return nil, app_errors.NotFoundError("CLIENT_PROJECTION_NOT_FOUND")
	}
	return projection, nil
}
