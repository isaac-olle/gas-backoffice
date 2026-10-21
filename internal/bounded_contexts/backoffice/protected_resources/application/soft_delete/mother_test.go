package soft_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/stretchr/testify/require"
)

// aPR builds a valid ProtectedResource owned by userId, ready for mutation (creation event drained).
func aPR(t *testing.T, userId ids2.UserId) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(ids2.NewUserId(), userId)
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "Resource", ownership, false)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}
