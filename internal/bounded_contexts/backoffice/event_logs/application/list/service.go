package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Limit caps how many of the most recent events are returned.
const Limit = 100

// Service lists the most recent events recorded by the audit subscriber.
type Service struct {
	repository event_logs.Repository
}

func NewService(repository event_logs.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, _ *Query) (any, error) {
	result, err := s.repository.Latest(ctx, Limit)
	if err != nil {
		return nil, app_errors.InternalError(event_logs.ErrEventLogsNotListed).WithCause(err, false)
	}
	return result, nil
}
