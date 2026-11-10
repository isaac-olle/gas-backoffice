package enable_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/enable"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId string
	author  *vo.UserAuthor
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withAuthor(author *vo.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
}

func aCommand(t *testing.T, opts ...commandOpt) *enable.Command {
	t.Helper()
	adminId := ids2.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	b := &commandBuilder{
		scopeId: uuid.New().String(),
		author:  author,
	}
	for _, o := range opts {
		o(b)
	}
	return enable.NewCommand(b.scopeId, b.author)
}

// aDisabledScope builds a non-system, disabled Scope aggregate owned by the given protected resource/subscription.
func aDisabledScope(t *testing.T, scopeId ids2.ScopeId, prId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(prId, subscriptionId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.Disable()
	scope.PullDomainEvents()
	return scope
}

func authorFor(t *testing.T, subscriptionId ids2.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}
