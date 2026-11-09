package hard_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId             ids2.ScopeId
	protectedResourceId ids2.ProtectedResourceId
	userId              ids2.UserId
}

func withScopeId(scopeId ids2.ScopeId) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withProtectedResourceId(prId ids2.ProtectedResourceId) commandOpt {
	return func(b *commandBuilder) { b.protectedResourceId = prId }
}

func withUserId(userId ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	b := &commandBuilder{
		scopeId:             ids2.NewScopeId(),
		protectedResourceId: ids2.NewProtectedResourceId(),
		userId:              ids2.NewUserId(),
	}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.scopeId, b.protectedResourceId, ids2.NewUserId(), b.userId)
}

// aScope builds a scope owned by a random protected resource/user, returning
// its id alongside the ownership ids so a matching Command can be built.
func aScope(t *testing.T) (*scopes.Scope, ids2.ScopeId, ids2.ProtectedResourceId, ids2.UserId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	prId := ids2.NewProtectedResourceId()
	userId := ids2.NewUserId()
	ownership := scopes.NewOwnershipFromVos(prId, ids2.NewUserId(), userId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope, scopeId, prId, userId
}
