package get

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Service struct {
	scopeFinder *domain_services.ScopeFinder
}

func NewService(scopeFinder *domain_services.ScopeFinder) *Service {
	return &Service{scopeFinder: scopeFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	scope, err := s.scopeFinder.Find(ctx, qry.scopeId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	subscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if err = scope.EnsureOwnership(qry.protectedResourceId, []ids.SubscriptionId{subscriptionId}); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return scope, nil
}
