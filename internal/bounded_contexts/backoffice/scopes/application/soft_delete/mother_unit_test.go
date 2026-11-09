package soft_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/soft_delete"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	scopeId             string
	protectedResourceId string
	userId              string
}

func withScopeId(scopeId string) commandOpt {
	return func(b *commandBuilder) { b.scopeId = scopeId }
}

func withProtectedResourceId(protectedResourceId string) commandOpt {
	return func(b *commandBuilder) { b.protectedResourceId = protectedResourceId }
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

// aCommand builds a soft_delete.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{
		scopeId:             uuid.New().String(),
		protectedResourceId: ids2.NewProtectedResourceId().String(),
		userId:              ids2.NewUserId().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return soft_delete.NewCommand(b.scopeId, b.protectedResourceId, ids2.NewUserId().String(), b.userId)
}

// aScope builds a non-system Scope aggregate owned by the given protected resource/user.
func aScope(t *testing.T, scopeId ids2.ScopeId, prId ids2.ProtectedResourceId, userId ids2.UserId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(prId, ids2.NewUserId(), userId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
