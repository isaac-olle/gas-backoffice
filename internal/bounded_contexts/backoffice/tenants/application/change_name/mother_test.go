package change_name_test

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/change_name"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	tenantId ids2.TenantId
	name     string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withTenantId(id ids2.TenantId) commandOpt {
	return func(b *commandBuilder) { b.tenantId = id }
}

// aCommand builds a change_name.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *change_name.Command {
	t.Helper()
	b := &commandBuilder{
		tenantId: ids2.NewTenantId(),
		name:     "New Name",
	}
	for _, o := range opts {
		o(b)
	}
	return change_name.NewCommand(b.tenantId, b.name)
}

// aTenant builds a Tenant aggregate bound to the given tenantId and adminId.
func aTenant(t *testing.T, tenantId ids2.TenantId, adminId ids2.AdminId, name string) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(tenantId, name, tenants.NewOwnership(adminId))
	require.NoError(t, err)
	tenant.PullDomainEvents() // clear creation event
	return tenant
}
