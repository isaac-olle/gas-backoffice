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

	allowedSubscriptions := make([]ids.SubscriptionId, 0, len(qry.allowedSubscriptions))
	for _, raw := range qry.allowedSubscriptions {
		subscriptionId, parseErr := ids.NewSubscriptionIdFromString(raw)
		if parseErr != nil {
			return nil, app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
				WithDescription("invalid allowed subscription id").WithCause(parseErr, false)
		}
		allowedSubscriptions = append(allowedSubscriptions, subscriptionId)
	}
	if len(qry.allowedSubscriptions) == 0 && !qry.allowAll {
		if subscriptionId, ok := ids.SubscriptionIdFromContext(ctx); ok {
			allowedSubscriptions = append(allowedSubscriptions, subscriptionId)
		}
	}
	if !qry.allowAll {
		if err = pr.EnsureOwnership(allowedSubscriptions); err != nil {
			return nil, app_errors.AddTrace(err)
		}
	}

	return pr.ToPrimitives().SubscriptionId, nil
}
