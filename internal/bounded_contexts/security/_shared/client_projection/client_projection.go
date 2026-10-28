package client_projection

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ClientProjection struct {
	ClientId       string
	SubscriptionId string
}

// EnsureSubscriptionAccess authorizes a client projection without loading the
// complete Client aggregate.
func (p *ClientProjection) EnsureSubscriptionAccess(allowed []ids.SubscriptionId) error {
	for _, id := range allowed {
		if id.String() == p.SubscriptionId {
			return nil
		}
	}
	return app_errors.ForbiddenError("SUBSCRIPTION_ACCESS_DENIED").
		WithDescription("client does not belong to an accessible subscription")
}

// IsOwnedBy allows the caller when it is the client's owning user, or when it
// is the client's owning admin. An admin is its own admin (callerAdminId ==
// callerUserId), which is what distinguishes it from a tenant: a tenant's
// AdminId is its owning admin, so without that check any tenant could reach the
// clients of every sibling tenant under the same admin.
