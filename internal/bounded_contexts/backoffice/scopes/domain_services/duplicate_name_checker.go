package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
)

type DuplicateNameChecker struct {
	repository scopes.Repository
}

func NewDuplicateNameChecker(repository scopes.Repository) *DuplicateNameChecker {
	return &DuplicateNameChecker{repository: repository}
}

func (c *DuplicateNameChecker) Check(ctx context.Context, name, protectedResourceId string) error {
	results, err := c.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: "name", Operator: criteria.Equals, Value: name},
			{Field: "protected_resource_id", Operator: criteria.Equals, Value: protectedResourceId},
		},
	})
	if err != nil {
		return app_errors.InternalError(app_errors.ErrQueryingDatabase).
			WithDescription("failed to query existing scopes").WithCause(err, false)
	}
	if len(results) > 0 {
		return app_errors.ConflictError(scopes.ErrScopeDuplicatedName).
			WithDescription("a scope with this name already exists")
	}
	return nil
}
