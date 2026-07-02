package get_test

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/stretchr/testify/require"
)

func aTenantWithId(t *testing.T, tenantId ids2.TenantId) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(tenantId, "Test Tenant", tenants.NewOwnership(ids2.NewAdminId()))
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}
