package hard_delete_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

// aPR builds a valid ProtectedResource ready for mutation (creation event drained).
func aPR(t *testing.T) *protected_resources.ProtectedResource {
	t.Helper()
	return aPROwnedBy(t, ids2.NewSubscriptionId())
}

// aPROwnedBy builds a valid ProtectedResource owned by subscriptionId (creation event drained).
func aPROwnedBy(t *testing.T, subscriptionId ids2.SubscriptionId) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(subscriptionId)
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "Resource", ownership, false)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}
