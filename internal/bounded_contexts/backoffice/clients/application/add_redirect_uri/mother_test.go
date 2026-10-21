package add_redirect_uri_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId    ids2.ClientId
	adminId     ids2.UserId
	userId      ids2.UserId
	redirectUri string
}

func withRedirectUri(uri string) commandOpt {
	return func(b *commandBuilder) { b.redirectUri = uri }
}

func withUserId(userId ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *add_redirect_uri.Command {
	t.Helper()
	b := &commandBuilder{
		clientId:    ids2.NewClientId(),
		adminId:     ids2.NewUserId(),
		userId:      ids2.NewUserId(),
		redirectUri: "https://example.com/callback",
	}
	for _, o := range opts {
		o(b)
	}
	return add_redirect_uri.NewCommand(b.clientId, b.adminId, b.userId, b.redirectUri)
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
