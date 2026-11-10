package get_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

// aPR builds a valid ProtectedResource with the given id, owned by subscriptionId, and no scope ids.
func aPR(t *testing.T, prId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(subscriptionId)
	pr, err := protected_resources.NewProtectedResource(prId, "My Resource", ownership, false)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}
