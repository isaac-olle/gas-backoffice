package domain_services

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const SubscriptionAccessDenied app_errors.ErrorCode = "SUBSCRIPTION_ACCESS_DENIED"

// SubscriptionAuthorizer checks access to an already-loaded resource's
// subscription. It does not perform repository lookups.
type SubscriptionAuthorizer struct{}

type SubscriptionOwned interface {
	EnsureSubscriptionAccess(allowedSubscriptions []ids.SubscriptionId) error
}

func NewSubscriptionAuthorizer() *SubscriptionAuthorizer { return &SubscriptionAuthorizer{} }

func (SubscriptionAuthorizer) Check(resource SubscriptionOwned, author *vo.UserAuthor) error {
	if author == nil {
		return app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	if author.Role.IsSystemAdministrator() {
		return nil
	}
	return resource.EnsureSubscriptionAccess(author.SubscriptionIDs())
}

func (SubscriptionAuthorizer) AllowedSubscriptionIDs(author *vo.UserAuthor) ([]string, bool, error) {
	if author == nil {
		return nil, false, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	if author.Role.IsSystemAdministrator() {
		return nil, true, nil
	}
	allowed := make([]string, 0, len(author.Memberships))
	for _, subscriptionID := range author.SubscriptionIDs() {
		allowed = append(allowed, subscriptionID.String())
	}
	return allowed, false, nil
}
