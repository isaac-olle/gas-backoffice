package list_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/list"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("filters a tenant caller by user_id and returns the matches", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		callerId := ids.NewUserId()
		expected := []*protected_resources.ProtectedResource{}
		want := criteria.Criteria{Filters: []criteria.Filter{{Field: "user_id", Operator: criteria.Equals, Value: callerId.String()}}}

		repo.EXPECT().Matching(mock.Anything, want).Return(expected, nil)

		result, err := list.NewService(repo).List(context.Background(), list.NewQuery(callerId, ids.NewUserId(), vo.TenantRole))
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("filters an admin caller by admin_id and returns the matches", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		adminId := ids.NewUserId()
		expected := []*protected_resources.ProtectedResource{}
		want := criteria.Criteria{Filters: []criteria.Filter{{Field: "admin_id", Operator: criteria.Equals, Value: adminId.String()}}}

		repo.EXPECT().Matching(mock.Anything, want).Return(expected, nil)

		result, err := list.NewService(repo).List(context.Background(), list.NewQuery(adminId, adminId, vo.AdminRole))
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		_, err := list.NewService(repo).List(context.Background(), list.NewQuery(ids.NewUserId(), ids.NewUserId(), vo.TenantRole))
		assert.Error(t, err)
	})
}
