package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/stretchr/testify/require"
)

func aTenantWithId(t *testing.T, userId ids2.UserId) *tenants.Tenant {
	t.Helper()
	return aTenantWithIdAndAdmin(t, userId, ids2.NewUserId())
}

func aTenantWithIdAndAdmin(t *testing.T, userId, adminId ids2.UserId) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(userId, "Test User", adminId, "", false)
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}
