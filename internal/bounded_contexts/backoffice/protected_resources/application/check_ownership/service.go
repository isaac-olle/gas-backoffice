package check_ownership

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Service struct {
	prFinder *domain_services.ProtectedResourceFinder
}

func NewService(prFinder *domain_services.ProtectedResourceFinder) *Service {
	return &Service{prFinder: prFinder}
}

// Check returns protected_resources.OwnershipPrimitives on success, or a
// app_errors.ForbiddenError (via EnsureOwnership) on failure.
func (s *Service) Check(ctx context.Context, qry *Query) (any, error) {
	prId, err := ids.NewProtectedResourceIdFromString(qry.protectedResourceId)
	if err != nil {
		return nil, app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("invalid protected resource id").WithCause(err, false)
	}

	pr, err := s.prFinder.Find(ctx, prId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	subscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if err = pr.EnsureOwnership([]ids.SubscriptionId{subscriptionId}); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return pr.ToPrimitives().SubscriptionId, nil
}
