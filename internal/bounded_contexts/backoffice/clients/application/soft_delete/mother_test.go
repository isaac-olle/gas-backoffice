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
	adminId  string
	userId   string
}

func withClientId(id string) commandOpt {
	return func(b *commandBuilder) { b.clientId = id }
}

func withUserId(id string) commandOpt {
	return func(b *commandBuilder) { b.userId = id }
}

func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{
		clientId: ids2.NewClientId().String(),
		adminId:  ids2.NewUserId().String(),
		userId:   ids2.NewUserId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return soft_delete.NewCommand(b.clientId, b.adminId, b.userId)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClient(t *testing.T, userId ids2.UserId) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(ids2.NewUserId(), userId),
		nil,
		nil,
		false,
		false,
	)
	require.NoError(t, err)
	c.PullDomainEvents() // clear creation event
	return c, clientId
}
