package remove_allowed_client_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

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

func aCommand(t *testing.T, opts ...commandOpt) *remove_allowed_client.Command {
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
	return remove_allowed_client.NewCommand(b.scopeId, b.clientId, b.author)
}

func aScopeWithClient(t *testing.T, clientId ids2.ClientId, prId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) (*scopes.Scope, ids2.ScopeId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	prims := scopes.ScopePrimitives{
		Id:                  scopeId.String(),
		Name:                "read:users",
		ProtectedResourceId: prId.String(),
		SubscriptionId:      subscriptionId.String(),
		AllowedClientIds:    []string{clientId.String()},
	}
	scope, err := scopes.FromPrimitives(prims)
	require.NoError(t, err)
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
