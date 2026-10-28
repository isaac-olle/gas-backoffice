package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

type Service struct {
	protectedResourceFinder *domain_services.ProtectedResourceFinder
	authorizer              *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(protectedResourceFinder *domain_services.ProtectedResourceFinder, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{protectedResourceFinder: protectedResourceFinder, authorizer: authorizer}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	pr, err := s.protectedResourceFinder.Find(ctx, qry.protectedResourceId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	if err = s.authorizer.Check(pr, qry.author); err != nil {
		return nil, app_errors.AddTrace(err)
	}
	return pr, nil
}
