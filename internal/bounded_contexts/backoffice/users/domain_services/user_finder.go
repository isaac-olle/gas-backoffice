package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

type UserFinder struct {
	repository users.Repository
}

func NewUserFinder(repository users.Repository) *UserFinder {
	return &UserFinder{repository: repository}
}

func (f *UserFinder) Find(ctx context.Context, id ids.UserId) (*users.User, error) {
	user, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(users.ErrUserNotFound).
			WithDescription("error searching user").WithCause(

			err, false)

	}
	if user == nil {
		return nil, app_errors.NotFoundError(users.ErrUserNotFound).
			WithDescription("user not found")
	}
	return user, nil
}
