package remove_allowed_grant_type_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_grant_type"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	clientId  ids2.ClientId
	adminId   ids2.UserId
	userId    ids2.UserId
	grantType string
}

func withGrantType(grantType string) commandOpt {
	return func(b *commandBuilder) { b.grantType = grantType }
}

func withUserId(userId ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *remove_allowed_grant_type.Command {
	t.Helper()
	b := &commandBuilder{
		clientId:  ids2.NewClientId(),
		adminId:   ids2.NewUserId(),
		userId:    ids2.NewUserId(),
		grantType: "client_credentials",
	}
	for _, o := range opts {
		o(b)
	}
	return remove_allowed_grant_type.NewCommand(b.clientId, b.adminId, b.userId, b.grantType)
}

func aClientWithGrantType(t *testing.T, userId ids2.UserId, grantType string) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		"valid-client",
		clients.NewOwnership(ids2.NewUserId(), userId),
		nil,
		[]string{grantType},
		false,
		false,
	)
	require.NoError(t, err)
	c.PullDomainEvents()
	return c, clientId
}
