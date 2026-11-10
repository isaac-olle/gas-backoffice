package get_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// ── Query Mother ───────────────────────────────────────────────────────────────

type queryOpt func(*queryBuilder)

type queryBuilder struct {
	clientId string
	author   *vo.UserAuthor
}

func withClientId(id string) queryOpt {
	return func(b *queryBuilder) { b.clientId = id }
}

func withAuthor(author *vo.UserAuthor) queryOpt {
	return func(b *queryBuilder) { b.author = author }
}

func aQuery(t *testing.T, opts ...queryOpt) *get.Query {
	t.Helper()
	adminId := ids2.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	b := &queryBuilder{
		clientId: uuid.New().String(),
		author:   author,
	}
	for _, o := range opts {
		o(b)
	}
	return get.NewQuery(b.clientId, b.author)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClientWithId(t *testing.T, name string, subscriptionId ids2.SubscriptionId) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		name,
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
