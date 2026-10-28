package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
)

// Service lists general users created by the caller.
type Service struct {
	repository users.Repository
}

func NewService(repository users.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, qry *Query) (any, error) {
	result, err := s.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{{Field: "author_id", Operator: criteria.Equals, Value: qry.adminId.String()}},
	})
	if err != nil {
		return nil, app_errors.InternalError(users.ErrUsersNotListed).WithCause(err, false)
	}
	return result, nil
}
