package soft_delete_test

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/soft_delete"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	tenantId string
}

func withTenantId(tenantId string) commandOpt {
	return func(b *commandBuilder) { b.tenantId = tenantId }
}

// aCommand builds a soft_delete.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{tenantId: uuid.New().String()}
	for _, o := range opts {
		o(b)
	}
	id, err := ids2.NewTenantIdFromString(b.tenantId)
	require.NoError(t, err)
	return soft_delete.NewCommand(id)
}

// aTenant builds a Tenant aggregate for use in tests.
func aTenant(t *testing.T) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(ids2.NewTenantId(), "Test Tenant", tenants.NewOwnership(ids2.NewAdminId()))
	require.NoError(t, err)
	tenant.PullDomainEvents() // clear creation event
	return tenant
}

// aTenantWithId builds a Tenant aggregate bound to the given tenantId.
func aTenantWithId(t *testing.T, tenantId ids2.TenantId) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(tenantId, "Test Tenant", tenants.NewOwnership(ids2.NewAdminId()))
	require.NoError(t, err)
	tenant.PullDomainEvents() // clear creation event
	return tenant
}
