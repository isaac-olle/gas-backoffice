package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/stretchr/testify/require"
)

// anAdmin builds a valid Admin aggregate from primitives for test setup.
func anAdmin(t *testing.T, adminId ids2.AdminId) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:        adminId.String(),
		Name:      "Test",
		TenantIds: []string{},
	})
	require.NoError(t, err)
	return a
}

// anAdminWithTenant builds an Admin that already owns one tenant.
func anAdminWithTenant(t *testing.T, adminId ids2.AdminId, tenantId ids2.TenantId) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:        adminId.String(),
		Name:      "Test",
		TenantIds: []string{tenantId.String()},
	})
	require.NoError(t, err)
	return a
}
