package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Limit caps how many of the most recent errors are returned.
const Limit = 100

// Service lists the most recent errors persisted by the error manager.
type Service struct {
	repository error_logs.Repository
}

func NewService(repository error_logs.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, _ *Query) (any, error) {
	result, err := s.repository.Latest(ctx, Limit)
	if err != nil {
		return nil, app_errors.InternalError(error_logs.ErrErrorLogsNotListed).WithCause(err, false)
	}
	return result, nil
}
