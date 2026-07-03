package soft_delete

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

func (s *Service) SoftDelete(ctx context.Context, cmd *Command) error {
	scopeId, err := ids.NewScopeIdFromString(cmd.scopeId)
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid scope id").
			WithCause(err)
	}

	if err = s.repository.SoftDelete(ctx, scopeId); err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to soft delete scope").
			WithCause(err)
	}
	return nil
}
