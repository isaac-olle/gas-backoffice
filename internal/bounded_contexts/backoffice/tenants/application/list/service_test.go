package list_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/list"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("filters by admin_id and returns the matches", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		callerId := ids.NewUserId()
		expected := []*tenants.Tenant{}
		want := criteria.Criteria{Filters: []criteria.Filter{{Field: "admin_id", Operator: criteria.Equals, Value: callerId.String()}}}

		repo.EXPECT().Matching(mock.Anything, want).Return(expected, nil)

		result, err := list.NewService(repo).List(context.Background(), list.NewQuery(callerId))
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		_, err := list.NewService(repo).List(context.Background(), list.NewQuery(ids.NewUserId()))
		assert.Error(t, err)
	})
}
