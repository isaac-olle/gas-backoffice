package add_allowed_client_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId             string
	clientId            string
	protectedResourceId string
	userId              string
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withClientId(clientId string) commandOpt {
	return func(b *commandBuilder) { b.clientId = clientId }
}

func withProtectedResourceId(protectedResourceId string) commandOpt {
	return func(b *commandBuilder) { b.protectedResourceId = protectedResourceId }
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

// aCommand builds an add_allowed_client.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *add_allowed_client.Command {
	t.Helper()
	b := &commandBuilder{
		scopeId:             ids2.NewScopeId().String(),
		clientId:            ids2.NewClientId().String(),
		protectedResourceId: ids2.NewProtectedResourceId().String(),
		userId:              ids2.NewUserId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return add_allowed_client.NewCommand(b.scopeId, b.clientId, b.protectedResourceId, ids2.NewUserId().String(), b.userId)
}

// aScope builds a Scope aggregate owned by the given protected resource/user for use in tests.
func aScope(t *testing.T, prId ids2.ProtectedResourceId, userId ids2.UserId) (*scopes.Scope, ids2.ScopeId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	ownership := scopes.NewOwnershipFromVos(prId, ids2.NewUserId(), userId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents() // clear creation event
	return scope, scopeId
}
