package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
)

// Service lists the tenants whose admin_id matches the caller.
type Service struct {
	repository tenants.Repository
}

func NewService(repository tenants.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, qry *Query) (any, error) {
	result, err := s.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{{Field: "admin_id", Operator: criteria.Equals, Value: qry.adminId.String()}},
	})
	if err != nil {
		return nil, app_errors.InternalError(tenants.ErrUsersNotListed).WithCause(err, false)
	}
	return result, nil
}
