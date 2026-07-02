package domain_services

import (
	"context"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type AdminFinder struct {
	repository admins.Repository
}

func NewAdminFinder(repository admins.Repository) *AdminFinder {
	return &AdminFinder{repository: repository}
}

func (f *AdminFinder) Find(ctx context.Context, id ids.AdminId) (*admins.Admin, error) {
	admin, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("error querying database for admin creation").
			WithCause(err)
	}
	if admin == nil {
		return nil, app_errors.NotFoundError(admins.ErrAdminNotFound).
			WithDescription("admin not  in database").
			WithDetails(map[string]interface{}{"admin_id": id.String()})
	}
	return admin, nil
}
