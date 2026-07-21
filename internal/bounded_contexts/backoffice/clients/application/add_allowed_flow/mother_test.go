package add_allowed_flow_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_flow"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId ids2.ClientId
	flow     string
}

func withFlow(flow string) commandOpt {
	return func(b *commandBuilder) { b.flow = flow }
}

func aCommand(t *testing.T, opts ...commandOpt) *add_allowed_flow.Command {
	t.Helper()
	b := &commandBuilder{
		clientId: ids2.NewClientId(),
		flow:     "pkce",
	}
	for _, o := range opts {
		o(b)
	}
	return add_allowed_flow.NewCommand(b.clientId, b.flow)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClient(t *testing.T) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(ids2.NewUserId()),
		nil,
		nil,
	)
	require.NoError(t, err)
	c.PullDomainEvents() // clear creation event
	return c, clientId
}
