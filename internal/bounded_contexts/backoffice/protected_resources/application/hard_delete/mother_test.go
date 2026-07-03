package hard_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/stretchr/testify/require"
)

// aPR builds a valid ProtectedResource ready for mutation (creation event drained).
func aPR(t *testing.T) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(ids2.NewTenantId(), ids2.NewAdminId())
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "Resource", ownership)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}
