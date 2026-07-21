package check_ownership_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
)

// ── Query Mother ───────────────────────────────────────────────────────────────

type queryOpt func(*queryBuilder)

type queryBuilder struct {
	clientId string
	userId   ids2.UserId
}

func withClientId(id string) queryOpt {
	return func(b *queryBuilder) { b.clientId = id }
}

func withUserId(id ids2.UserId) queryOpt {
	return func(b *queryBuilder) { b.userId = id }
}

func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	b := &queryBuilder{
		clientId: ids2.NewClientId().String(),
		userId:   ids2.NewUserId(),
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.clientId, b.userId.String())
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClientWithOwnership(t *testing.T, userId ids2.UserId) *clients.Client {
	t.Helper()
	c, err := clients.NewClient(
		ids2.NewClientId(),
		"valid-client",
		clients.NewOwnership(userId),
		nil,
		nil,
	)
	require.NoError(t, err)
	c.PullDomainEvents()
	return c
}
