package check_ownership_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/stretchr/testify/require"
)

// aPRWithOwnership builds a ProtectedResource with the given ids as ownership.
func aPRWithOwnership(
	t *testing.T,
	prId ids2.ProtectedResourceId,
	userId ids2.UserId,
	adminId ids2.AdminId,
) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(userId, adminId)
	pr, err := protected_resources.NewProtectedResource(prId, "R", ownership)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}
