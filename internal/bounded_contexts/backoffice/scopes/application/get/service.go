package get

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

type Service struct {
	scopeFinder *domain_services.ScopeFinder
	authorizer  *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(scopeFinder *domain_services.ScopeFinder, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{scopeFinder: scopeFinder, authorizer: authorizer}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	scope, err := s.scopeFinder.Find(ctx, qry.scopeId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	if err = s.authorizer.Check(scope, qry.author); err != nil {
		return nil, app_errors.AddTrace(err)
	}
	return scope, nil
}
