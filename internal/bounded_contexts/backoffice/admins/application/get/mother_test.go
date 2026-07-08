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
		Id:      adminId.String(),
		Name:    "Test",
		UserIds: []string{},
	})
	require.NoError(t, err)
	return a
}

// anAdminWithUser builds an Admin that already owns one user.
func anAdminWithUser(t *testing.T, adminId ids2.AdminId, userId ids2.UserId) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:      adminId.String(),
		Name:    "Test",
		UserIds: []string{userId.String()},
	})
	require.NoError(t, err)
	return a
}
