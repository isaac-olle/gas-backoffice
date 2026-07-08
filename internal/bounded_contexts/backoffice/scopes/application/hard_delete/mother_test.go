package hard_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId ids2.ScopeId
}

func withScopeId(scopeId ids2.ScopeId) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

// aCommand builds a hard_delete.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	b := &commandBuilder{scopeId: ids2.NewScopeId()}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.scopeId)
}

// aScope builds a Scope aggregate for use in tests.
func aScope(t *testing.T) (*scopes.Scope, ids2.ScopeId) {
	t.Helper()
	scopeId := ids2.NewScopeId()
	ownership := scopes.NewOwnershipFromVos(
		ids2.NewProtectedResourceId(),
		ids2.NewUserId(),
		ids2.NewAdminId(),
	)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership)
	require.NoError(t, err)
	scope.PullDomainEvents() // clear creation event
	return scope, scopeId
}
