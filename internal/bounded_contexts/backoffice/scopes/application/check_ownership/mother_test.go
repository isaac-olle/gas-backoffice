package check_ownership_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/check_ownership"
	"github.com/stretchr/testify/require"
)

type queryOpt func(*queryBuilder)

type queryBuilder struct {
	scopeId             string
	protectedResourceId string
	userId              string
}

func withScopeId(scopeId string) queryOpt {
	return func(b *queryBuilder) { b.scopeId = scopeId }
}

func withProtectedResourceId(prId string) queryOpt {
	return func(b *queryBuilder) { b.protectedResourceId = prId }
}

func withUserId(userId string) queryOpt {
	return func(b *queryBuilder) { b.userId = userId }
}

func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	b := &queryBuilder{
		scopeId:             ids2.NewScopeId().String(),
		protectedResourceId: ids2.NewProtectedResourceId().String(),
		userId:              ids2.NewUserId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.scopeId, b.protectedResourceId, b.userId)
}

func aScope(t *testing.T, scopeId ids2.ScopeId, prId ids2.ProtectedResourceId, userId ids2.UserId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(prId, userId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
