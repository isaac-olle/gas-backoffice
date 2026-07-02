package get_test

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
	"github.com/stretchr/testify/require"
)

func aScope(t *testing.T, scopeId ids2.ScopeId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewTenantId(), ids2.NewAdminId())
	scope, err := scopes.NewScope(scopeId, "read:users", ownership)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
