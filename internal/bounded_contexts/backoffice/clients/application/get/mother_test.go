package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/get"
)

// ── Query Mother ───────────────────────────────────────────────────────────────

type queryOpt func(*queryBuilder)

type queryBuilder struct {
	clientId string
	adminId  string
	userId   string
}

func withClientId(id string) queryOpt {
	return func(b *queryBuilder) { b.clientId = id }
}

func withUserId(id string) queryOpt {
	return func(b *queryBuilder) { b.userId = id }
}

func aQuery(t *testing.T, opts ...queryOpt) *get.Query {
	t.Helper()
	b := &queryBuilder{
		clientId: uuid.New().String(),
		adminId:  uuid.New().String(),
		userId:   uuid.New().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return get.NewQuery(b.clientId, b.adminId, b.userId)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClientWithId(t *testing.T, name string, userId ids2.UserId) (*clients.Client, ids2.ClientId) {
	t.Helper()
	clientId := ids2.NewClientId()
	c, err := clients.NewClient(
		clientId,
		name,
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
