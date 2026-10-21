package create_tenant_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/create_tenant"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	id       ids2.UserId
	adminId  ids2.UserId
	name     string
	isSystem bool
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withIsSystem(isSystem bool) commandOpt {
	return func(b *commandBuilder) { b.isSystem = isSystem }
}

func aCommand(t *testing.T, opts ...commandOpt) *create_tenant.Command {
	t.Helper()
	b := &commandBuilder{
		id:      ids2.NewUserId(),
		adminId: ids2.NewUserId(),
		name:    "Tenant Name",
	}
	for _, o := range opts {
		o(b)
	}
	cmd := create_tenant.NewCommand(b.id, b.adminId, b.name, b.isSystem)
	require.NotNil(t, cmd)
	return cmd
}
