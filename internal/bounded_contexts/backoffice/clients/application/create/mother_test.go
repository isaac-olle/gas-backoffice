package create_test

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/application/create"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	name         string
	adminId      ids2.AdminId
	tenantId     ids2.TenantId
	allowedFlows []string
	redirectUris []string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withFlows(flows ...string) commandOpt {
	return func(b *commandBuilder) { b.allowedFlows = flows }
}

func withRedirectUris(uris ...string) commandOpt {
	return func(b *commandBuilder) { b.redirectUris = uris }
}

func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{
		name:     "valid-client",
		adminId:  ids2.NewAdminId(),
		tenantId: ids2.NewTenantId(),
	}
	for _, o := range opts {
		o(b)
	}
	clientId := ids2.NewClientId()
	ownership := clients.NewOwnership(b.adminId, b.tenantId)
	cmd := create.NewCommand(clientId, b.name, ownership, b.allowedFlows, b.redirectUris)
	require.NotNil(t, cmd)
	return cmd
}
