package soft_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/soft_delete"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId string
}

func withClientId(id string) commandOpt {
	return func(b *commandBuilder) { b.clientId = id }
}

func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{
		clientId: ids2.NewClientId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return soft_delete.NewCommand(b.clientId)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClient(t *testing.T) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(ids2.NewAdminId(), ids2.NewTenantId()),
		nil,
		nil,
	)
	require.NoError(t, err)
	c.PullDomainEvents() // clear creation event
	return c, clientId
}
