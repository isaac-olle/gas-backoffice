package add_allowed_client_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId  string
	clientId string
	author   *vo.UserAuthor
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withClientId(clientId string) commandOpt {
	return func(b *commandBuilder) { b.clientId = clientId }
}

func withAuthor(author *vo.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
}

// aCommand builds an add_allowed_client.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *add_allowed_client.Command {
	t.Helper()
	adminId := ids2.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminId.String(), adminId.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	b := &commandBuilder{
		scopeId:  ids2.NewScopeId().String(),
		clientId: ids2.NewClientId().String(),
		author:   author,
	}
	for _, o := range opts {
		o(b)
	}
	return add_allowed_client.NewCommand(b.scopeId, b.clientId, b.author)
}

// aScope builds a Scope aggregate owned by the given protected resource/subscription for use in tests.
func aScope(t *testing.T, prId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) (*scopes.Scope, ids2.ScopeId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	ownership := scopes.NewOwnershipFromVos(prId, subscriptionId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents() // clear creation event
	return scope, scopeId
}

func authorFor(t *testing.T, subscriptionId ids2.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}

// clientProjectionRepoStub is a minimal scopes.ClientProjectionRepository implementation for tests.
type clientProjectionRepoStub struct {
	projection *scopes.ClientProjection
	err        error
}

func (s *clientProjectionRepoStub) FindClientProjection(_ context.Context, _ ids2.ClientId) (*scopes.ClientProjection, error) {
	return s.projection, s.err
}
