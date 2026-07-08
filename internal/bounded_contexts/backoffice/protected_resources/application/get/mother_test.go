package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/stretchr/testify/require"
)

// aPR builds a valid ProtectedResource with the given id and no scope ids.
func aPR(t *testing.T, prId ids2.ProtectedResourceId) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(ids2.NewUserId(), ids2.NewAdminId())
	pr, err := protected_resources.NewProtectedResource(prId, "My Resource", ownership)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}
