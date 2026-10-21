package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
)

type AdminFinder struct {
	repository admins.Repository
}

func NewAdminFinder(repository admins.Repository) *AdminFinder {
	return &AdminFinder{repository: repository}
}

func (f *AdminFinder) Find(ctx context.Context, id ids.UserId) (*admins.Admin, error) {
	admin, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(admins.ErrAdminNotFound).
			WithDescription("error searching admin").WithCause(

			err, false)

	}
	if admin == nil {
		return nil, app_errors.NotFoundError(admins.ErrAdminNotFound).
			WithDescription("admin not found")
	}
	return admin, nil
}
