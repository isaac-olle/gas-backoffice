package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/stretchr/testify/require"
)

func aScope(t *testing.T, scopeId ids2.ScopeId, prId ids2.ProtectedResourceId, userId ids2.UserId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(prId, ids2.NewUserId(), userId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
