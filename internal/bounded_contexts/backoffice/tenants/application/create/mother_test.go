package create_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/create"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	name    string
	adminId string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withAdminId(adminId string) commandOpt {
	return func(b *commandBuilder) { b.adminId = adminId }
}

// aCommand builds a create.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{
		name:    "my-tenant",
		adminId: uuid.New().String(),
	}
	for _, o := range opts {
		o(b)
	}

	tenantId, err := ids2.NewTenantIdFromString(uuid.New().String())
	require.NoError(t, err)
	adminId, err := ids2.NewAdminIdFromString(b.adminId)
	require.NoError(t, err)

	return create.NewCommand(tenantId, b.name, tenants.NewOwnership(adminId))
}
