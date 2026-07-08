package domain_services_test

import (
	"context"
	"fmt"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestScopeFinder_Find(t *testing.T) {
	t.Parallel()

	t.Run("returns scope when found", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		expected := aScope(t)

		repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(expected, nil)

		finder := domain_services.NewScopeFinder(repo)
		got, err := finder.Find(context.Background(), ids2.NewScopeId())

		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)

		repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewScopeFinder(repo)
		_, err := finder.Find(context.Background(), ids2.NewScopeId())

		assert.Error(t, err)
	})
}

// --- helpers ---

func aScope(t *testing.T) *scopes.Scope {
	t.Helper()
	ownership := scopes.NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewUserId(), ids2.NewAdminId())
	scope, err := scopes.NewScope(ids2.NewScopeId(), "read:users", ownership)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}
