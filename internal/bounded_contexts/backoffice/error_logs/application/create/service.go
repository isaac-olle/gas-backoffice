package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Service persists an error log. There is nothing to validate: the error
// manager is the only caller and it already built every field from an
// *app_errors.AppError, so no VOs are involved.
type Service struct {
	repository error_logs.Repository
}

func NewService(repository error_logs.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	errorLog := error_logs.ErrorLog{
		AuthorId:    cmd.authorId,
		AuthorRole:  cmd.authorRole,
		ErrorId:     cmd.errorId,
		Type:        cmd.errorType,
		Code:        cmd.code,
		Description: cmd.description,
		Details:     cmd.details,
		Cause:       cmd.cause,
		StackTrace:  cmd.stackTrace,
		RequestId:   cmd.requestId,
		OccurredAt:  cmd.occurredAt,
	}

	if err := s.repository.Save(ctx, errorLog); err != nil {
		return app_errors.InternalError(error_logs.ErrErrorLogNotSaved).WithCause(err, false)
	}
	return nil
}
