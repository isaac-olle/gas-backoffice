package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

// Service lists clients visible to the caller's allowed subscriptions.
type Service struct {
	repository clients.Repository
	authorizer *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(repository clients.Repository, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
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
		return nil, app_errors.InternalError(clients.ErrClientsNotListed).WithCause(err, false)
	}
	return result, nil
}
