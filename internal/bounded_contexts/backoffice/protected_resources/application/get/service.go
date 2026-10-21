package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Service struct {
	protectedResourceFinder *domain_services.ProtectedResourceFinder
}

func NewService(protectedResourceFinder *domain_services.ProtectedResourceFinder) *Service {
	return &Service{protectedResourceFinder: protectedResourceFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	pr, err := s.protectedResourceFinder.Find(ctx, qry.protectedResourceId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	subscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if err = pr.EnsureOwnership([]ids.SubscriptionId{subscriptionId}); err != nil {
		return nil, app_errors.AddTrace(err)
	}
	return pr, nil
}
