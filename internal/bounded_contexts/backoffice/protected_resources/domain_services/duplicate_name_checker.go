package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
)

type DuplicateNameChecker struct {
	repository protected_resources.Repository
}

func NewDuplicateNameChecker(repository protected_resources.Repository) *DuplicateNameChecker {
	return &DuplicateNameChecker{repository: repository}
}

func (c *DuplicateNameChecker) Check(ctx context.Context, name, subscriptionId string) error {
	results, err := c.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: "name", Operator: criteria.Equals, Value: name},
			{Field: "subscription_id", Operator: criteria.Equals, Value: subscriptionId},
		},
	})
	if err != nil {
		return app_errors.InternalError(app_errors.ErrQueryingDatabase).
			WithDescription("failed to query existing protected resources").WithCause(err, false)
	}
	if len(results) > 0 {
		return app_errors.ConflictError(protected_resources.ErrProtectedResourceDuplicateName).
			WithDescription("a protected resource with this name already exists")
	}
	return nil
}
