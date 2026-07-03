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
	protectedResourceId string
	tenantId            string
	adminId             string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

// aOwnership builds a scopes.Ownership using the Object Mother pattern.
func aOwnership(t *testing.T) scopes.Ownership {
	t.Helper()
	prId, err := ids2.NewProtectedResourceIdFromString(uuid.New().String())
	require.NoError(t, err)
	tenantId, err := ids2.NewTenantIdFromString(uuid.New().String())
	require.NoError(t, err)
	adminId, err := ids2.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)
	return scopes.NewOwnershipFromVos(prId, tenantId, adminId)
}

// aCommand builds a create.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{name: "tenants:read"}
	for _, o := range opts {
		o(b)
	}
	scopeId, err := ids2.NewScopeIdFromString(uuid.New().String())
	require.NoError(t, err)
	return create.NewCommand(scopeId, b.name, aOwnership(t))
}
