package hard_delete_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId ids2.ScopeId
	author  *vo.UserAuthor
}

func withScopeId(scopeId ids2.ScopeId) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withAuthor(author *vo.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
}

func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	adminId := ids2.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	b := &commandBuilder{
		scopeId: ids2.NewScopeId(),
		author:  author,
	}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.scopeId, b.author)
}

// aScope builds a scope owned by a random protected resource/subscription, returning
// its id alongside the ownership ids so a matching Command can be built.
func aScope(t *testing.T) (*scopes.Scope, ids2.ScopeId, ids2.ProtectedResourceId, ids2.SubscriptionId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	prId := ids2.NewProtectedResourceId()
	subscriptionId := ids2.NewSubscriptionId()
	ownership := scopes.NewOwnershipFromVos(prId, subscriptionId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope, scopeId, prId, subscriptionId
}

func authorFor(t *testing.T, subscriptionId ids2.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}
