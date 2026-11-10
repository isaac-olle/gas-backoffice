package remove_allowed_grant_type_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_grant_type"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId  ids2.ClientId
	author    *vo.UserAuthor
	grantType string
}

func withGrantType(grantType string) commandOpt {
	return func(b *commandBuilder) { b.grantType = grantType }
}

func withAuthor(author *vo.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
}

func aCommand(t *testing.T, opts ...commandOpt) *remove_allowed_grant_type.Command {
	t.Helper()
	adminId := ids2.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	b := &commandBuilder{
		clientId:  ids2.NewClientId(),
		author:    author,
		grantType: "client_credentials",
	}
	for _, o := range opts {
		o(b)
	}
	return remove_allowed_grant_type.NewCommand(b.clientId, b.grantType, b.author)
}

func aClientWithGrantType(t *testing.T, subscriptionId ids2.SubscriptionId, grantType string) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(subscriptionId),
		nil,
		[]string{grantType},
		false,
		false,
	)
	require.NoError(t, err)
	c.PullDomainEvents()
	return c, clientId
}
