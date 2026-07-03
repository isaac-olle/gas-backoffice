package check_ownership_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/check_ownership"
	"github.com/stretchr/testify/require"
)

// queryOpt is a functional option for queryBuilder.
type queryOpt func(*queryBuilder)

type queryBuilder struct {
	scopeId             string
	protectedResourceId string
	tenantId            string
	adminId             string
}

func withScopeId(scopeId string) queryOpt {
	return func(b *queryBuilder) { b.scopeId = scopeId }
}

func withProtectedResourceId(prId string) queryOpt {
	return func(b *queryBuilder) { b.protectedResourceId = prId }
}

func withTenantId(tenantId string) queryOpt {
	return func(b *queryBuilder) { b.tenantId = tenantId }
}

func withAdminId(adminId string) queryOpt {
	return func(b *queryBuilder) { b.adminId = adminId }
}

// aQuery builds a check_ownership.Query using the Object Mother pattern.
func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	b := &queryBuilder{
		scopeId:             ids2.NewScopeId().String(),
		protectedResourceId: ids2.NewProtectedResourceId().String(),
		tenantId:            ids2.NewTenantId().String(),
		adminId:             ids2.NewAdminId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.scopeId, b.protectedResourceId, b.tenantId, b.adminId)
}

// aScope builds a Scope aggregate owned by the given ids.
func aScope(t *testing.T, scopeId ids2.ScopeId, prId ids2.ProtectedResourceId, tenantId ids2.TenantId, adminId ids2.AdminId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(prId, tenantId, adminId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
