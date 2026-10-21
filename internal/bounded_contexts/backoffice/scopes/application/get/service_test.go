package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(repo *scopemocks.MockRepository) *get.Service {
	return get.NewService(domain_services.NewScopeFinder(repo))
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns scope when found", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		scopeId := ids.NewScopeId()
		prId := ids.NewProtectedResourceId()
		userId := ids.NewUserId()
		scope := aScope(t, scopeId, prId, userId)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)

		svc := newService(repo)
		result, err := svc.Get(context.Background(), get.NewQuery(scopeId, prId, ids.NewUserId(), userId))
		require.NoError(t, err)

		s, ok := result.(*scopes.Scope)
		require.True(t, ok)
		assert.Equal(t, scopeId.String(), s.ToPrimitives().Id)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		scopeId := ids.NewScopeId()
		prId := ids.NewProtectedResourceId()
		userId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(nil, fmt.Errorf("db error"))

		svc := newService(repo)
		_, err := svc.Get(context.Background(), get.NewQuery(scopeId, prId, ids.NewUserId(), userId))
		assert.Error(t, err)
	})
}
