package add_redirect_uri_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

// ── Command Mother ─────────────────────────────────────────────────────────────

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId    ids2.ClientId
	author      *vo.UserAuthor
	redirectUri string
}

func withRedirectUri(uri string) commandOpt {
	return func(b *commandBuilder) { b.redirectUri = uri }
}

func withAuthor(author *vo.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
}

func aCommand(t *testing.T, opts ...commandOpt) *add_redirect_uri.Command {
	t.Helper()
	adminId := ids2.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	b := &commandBuilder{
		clientId:    ids2.NewClientId(),
		author:      author,
		redirectUri: "https://example.com/callback",
	}
	for _, o := range opts {
		o(b)
	}
	return add_redirect_uri.NewCommand(b.clientId, b.redirectUri, b.author)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClient(t *testing.T, subscriptionId ids2.SubscriptionId) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(subscriptionId),
		nil,
		nil,
		false,
		false,
	)
	require.NoError(t, err)
	c.PullDomainEvents() // clear creation event
	return c, clientId
}
