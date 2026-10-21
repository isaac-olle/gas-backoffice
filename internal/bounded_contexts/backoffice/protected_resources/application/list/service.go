package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
)

// Service lists the protected resources belonging to the subscription in the
// request path.
type Service struct {
	repository protected_resources.Repository
}

func NewService(repository protected_resources.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, qry *Query) (any, error) {
	result, err := s.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{{Field: "subscription_id", Operator: criteria.Equals, Value: qry.subscriptionId.String()}},
	})
	if err != nil {
		return nil, app_errors.InternalError(protected_resources.ErrProtectedResourcesNotListed).WithCause(err, false)
	}
	return result, nil
}
