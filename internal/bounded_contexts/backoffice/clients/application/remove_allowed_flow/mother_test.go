package remove_allowed_flow_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_flow"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId ids2.ClientId
	flow     string
}

func withFlow(flow string) commandOpt {
	return func(b *commandBuilder) { b.flow = flow }
}

func aCommand(t *testing.T, opts ...commandOpt) *remove_allowed_flow.Command {
	t.Helper()
	b := &commandBuilder{
		clientId: ids2.NewClientId(),
		flow:     "pkce",
	}
	for _, o := range opts {
		o(b)
	}
	return remove_allowed_flow.NewCommand(b.clientId, b.flow)
}

func aClientWithFlow(t *testing.T, flow string) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(ids2.NewAdminId(), ids2.NewUserId()),
		nil,
		[]string{flow},
	)
	require.NoError(t, err)
	c.PullDomainEvents()
	return c, clientId
}
