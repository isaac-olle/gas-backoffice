package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

// Service lists protected resources visible to the caller's allowed
// subscriptions.
type Service struct {
	repository protected_resources.Repository
	authorizer *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(repository protected_resources.Repository, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (s *Service) List(ctx context.Context, qry *Query) (any, error) {
	filterCriteria := criteria.Criteria{}
	allowed, allowAll, err := s.authorizer.AllowedSubscriptionIDs(qry.author)
	if err != nil {
		return nil, err
	}
	if !allowAll {
		filterCriteria.Filters = []criteria.Filter{{Field: "subscription_id", Operator: criteria.AnyOf, Value: allowed}}
	}
	result, err := s.repository.Matching(ctx, filterCriteria)
	if err != nil {
		return nil, app_errors.InternalError(protected_resources.ErrProtectedResourcesNotListed).WithCause(err, false)
	}
	return result, nil
}
