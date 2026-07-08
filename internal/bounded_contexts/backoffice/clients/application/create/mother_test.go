package create_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	name         string
	adminId      ids2.AdminId
	userId       ids2.UserId
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
		name:    "valid-client",
		adminId: ids2.NewAdminId(),
		userId:  ids2.NewUserId(),
	}
	for _, o := range opts {
		o(b)
	}
	clientId := ids2.NewClientId()
	ownership := clients.NewOwnership(b.adminId, b.userId)
	cmd := create.NewCommand(clientId, b.name, ownership, b.allowedFlows, b.redirectUris)
	require.NotNil(t, cmd)
	return cmd
}
