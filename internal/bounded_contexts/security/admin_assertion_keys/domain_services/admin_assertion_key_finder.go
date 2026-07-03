package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
)

type AdminAssertionKeyFinder struct {
	repository admin_assertion_keys.Repository
}

func NewAdminAssertionKeyFinder(repository admin_assertion_keys.Repository) *AdminAssertionKeyFinder {
	return &AdminAssertionKeyFinder{repository: repository}
}

func (f *AdminAssertionKeyFinder) Find(ctx context.Context, id ids.AdminAssertionKeyId) (*admin_assertion_keys.AdminAssertionKey, error) {
	aak, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(admin_assertion_keys.ErrAdminAssertionKeyQueryFailed).
			WithDescription("repository failed when sending the query to database").
			WithCause(err)
	}
	if aak == nil {
		return nil, app_errors.NotFoundError(admin_assertion_keys.ErrAdminAssertionKeyNotFound).
			WithDescription("admin assertion key not found").
			WithDetails(map[string]any{"key_id": id.String()})
	}
	return aak, nil
}
