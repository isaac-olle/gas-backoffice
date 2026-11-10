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
	name              string
	subscriptionId    ids2.SubscriptionId
	allowedGrantTypes []string
	redirectUris      []string
	public            bool
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withGrantTypes(grantTypes ...string) commandOpt {
	return func(b *commandBuilder) { b.allowedGrantTypes = grantTypes }
}

func withRedirectUris(uris ...string) commandOpt {
	return func(b *commandBuilder) { b.redirectUris = uris }
}

func withPublic(public bool) commandOpt {
	return func(b *commandBuilder) { b.public = public }
}

func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{
		name:           "valid-client",
		subscriptionId: ids2.NewSubscriptionId(),
	}
	for _, o := range opts {
		o(b)
	}
	clientId := ids2.NewClientId()
	ownership := clients.NewOwnership(b.subscriptionId)
	cmd := create.NewCommand(clientId, b.name, ownership, b.allowedGrantTypes, b.redirectUris, b.public, false)
	require.NotNil(t, cmd)
	return cmd
}
