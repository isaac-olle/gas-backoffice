package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/stretchr/testify/require"
)

func anAdminWithId(t *testing.T, adminId ids2.UserId) *admins.Admin {
	t.Helper()
	admin, err := admins.NewAdmin(adminId, "Test Admin", "", false)
	require.NoError(t, err)
	admin.PullDomainEvents()
	return admin
}
