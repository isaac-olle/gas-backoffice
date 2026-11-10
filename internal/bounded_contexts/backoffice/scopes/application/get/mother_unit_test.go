package get_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

func aScope(t *testing.T, scopeId ids2.ScopeId, prId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(prId, subscriptionId)
	scope, err := scopes.NewScope(scopeId, "read:users", ownership, false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
