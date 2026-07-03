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
	adminId  ids2.AdminId
	tenantId ids2.TenantId
}

func withClientId(id string) queryOpt {
	return func(b *queryBuilder) { b.clientId = id }
}

func withAdminId(id ids2.AdminId) queryOpt {
	return func(b *queryBuilder) { b.adminId = id }
}

func withTenantId(id ids2.TenantId) queryOpt {
	return func(b *queryBuilder) { b.tenantId = id }
}

func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	b := &queryBuilder{
		clientId: ids2.NewClientId().String(),
		adminId:  ids2.NewAdminId(),
		tenantId: ids2.NewTenantId(),
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.clientId, b.adminId, b.tenantId)
}

// ── Client Mother ──────────────────────────────────────────────────────────────

func aClientWithOwnership(t *testing.T, adminId ids2.AdminId, tenantId ids2.TenantId) *clients.Client {
	t.Helper()
	c, err := clients.NewClient(
		ids2.NewClientId(),
		"valid-client",
		clients.NewOwnership(adminId, tenantId),
		nil,
		nil,
	)
	require.NoError(t, err)
	c.PullDomainEvents()
	return c
}
