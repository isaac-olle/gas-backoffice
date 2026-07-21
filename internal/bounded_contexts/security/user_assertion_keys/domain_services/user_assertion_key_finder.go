package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys"
)

type UserAssertionKeyFinder struct {
	repository user_assertion_keys.Repository
}

func NewUserAssertionKeyFinder(repository user_assertion_keys.Repository) *UserAssertionKeyFinder {
	return &UserAssertionKeyFinder{repository: repository}
}

func (f *UserAssertionKeyFinder) Find(ctx context.Context, id ids.UserAssertionKeyId) (*user_assertion_keys.UserAssertionKey, error) {
	uak, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(user_assertion_keys.ErrUserAssertionKeyQueryFailed).
			WithDescription("repository failed when sending the query to database").
			WithCause(err)
	}
	if uak == nil {
		return nil, app_errors.NotFoundError(user_assertion_keys.ErrUserAssertionKeyNotFound).
			WithDescription("user assertion key not found").
			WithDetails(map[string]any{"key_id": id.String()})
	}
	return uak, nil
}
