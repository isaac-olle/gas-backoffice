package check_ownership_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/check_ownership"
	"github.com/stretchr/testify/require"
)

// queryOpt is a functional option for queryBuilder.
type queryOpt func(*queryBuilder)

type queryBuilder struct {
	tenantId string
	adminId  ids2.AdminId
}

func withTenantId(tenantId string) queryOpt {
	return func(b *queryBuilder) { b.tenantId = tenantId }
}

func withAdminId(adminId ids2.AdminId) queryOpt {
	return func(b *queryBuilder) { b.adminId = adminId }
}

// aQuery builds a check_ownership.Query using the Object Mother pattern.
func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	adminId, err := ids2.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)
	b := &queryBuilder{
		tenantId: uuid.New().String(),
		adminId:  adminId,
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.tenantId, b.adminId)
}

// aTenant builds a Tenant aggregate owned by the given admin.
func aTenant(t *testing.T, tenantId ids2.TenantId, adminId ids2.AdminId) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(tenantId, "Test Tenant", tenants.NewOwnership(adminId))
	require.NoError(t, err)
	tenant.PullDomainEvents() // clear creation event
	return tenant
}
