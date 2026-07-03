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
	redirectUri string
}

func withRedirectUri(uri string) commandOpt {
	return func(b *commandBuilder) { b.redirectUri = uri }
}

func aCommand(t *testing.T, opts ...commandOpt) *add_redirect_uri.Command {
	t.Helper()
	b := &commandBuilder{
		clientId:    ids2.NewClientId(),
		redirectUri: "https://example.com/callback",
	}
	for _, o := range opts {
		o(b)
	}
	return add_redirect_uri.NewCommand(b.clientId, b.redirectUri)
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
