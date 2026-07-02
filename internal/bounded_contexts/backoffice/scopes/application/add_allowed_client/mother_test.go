package add_allowed_client_test

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId  string
	clientId string
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withClientId(clientId string) commandOpt {
	return func(b *commandBuilder) { b.clientId = clientId }
}

// aCommand builds an add_allowed_client.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *add_allowed_client.Command {
	t.Helper()
	b := &commandBuilder{
		scopeId:  ids2.NewScopeId().String(),
		clientId: ids2.NewClientId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return add_allowed_client.NewCommand(b.scopeId, b.clientId)
}

// aScope builds a Scope aggregate for use in tests.
func aScope(t *testing.T) (*scopes.Scope, ids2.ScopeId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	ownership := scopes.NewOwnershipFromVos(
		ids2.NewProtectedResourceId(),
		ids2.NewTenantId(),
		ids2.NewAdminId(),
	)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership)
	require.NoError(t, err)
	scope.PullDomainEvents() // clear creation event
	return scope, scopeId
}
