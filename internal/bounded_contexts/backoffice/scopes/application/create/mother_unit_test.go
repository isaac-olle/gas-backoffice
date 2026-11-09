package create_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	name                string
	protectedResourceId ids2.ProtectedResourceId
	userId              ids2.UserId
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withProtectedResourceId(prId ids2.ProtectedResourceId) commandOpt {
	return func(b *commandBuilder) { b.protectedResourceId = prId }
}

func withUserId(userId ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

// aCommand builds a create.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{
		name:                "tenants:read",
		protectedResourceId: ids2.NewProtectedResourceId(),
		userId:              ids2.NewUserId(),
	}
	for _, o := range opts {
		o(b)
	}
	scopeId, err := ids2.NewScopeIdFromString(uuid.New().String())
	require.NoError(t, err)
	ownership := scopes.NewOwnershipFromVos(b.protectedResourceId, ids2.NewUserId(), b.userId)
	return create.NewCommand(scopeId, b.name, ownership, false)
}
