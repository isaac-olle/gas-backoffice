package get_test

import (
	"context"
	"fmt"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	adminmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*get.Service, *adminmocks.MockRepository) {
	t.Helper()
	repo := adminmocks.NewMockRepository(t)
	finder := domain_services.NewAdminFinder(repo)
	return get.NewService(finder), repo
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns admin primitives with correct id and name", func(t *testing.T) {
		t.Parallel()
		svc, repo := newService(t)
		adminId := ids2.NewAdminId()
		admin := anAdmin(t, adminId)

		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)

		result, err := svc.Get(context.Background(), get.NewQuery(adminId))
		require.NoError(t, err)

		prims, ok := result.(admins.AdminPrimitives)
		require.True(t, ok)
		assert.Equal(t, adminId.String(), prims.Id)
		assert.Equal(t, "Test", prims.Name)
		assert.Empty(t, prims.UserIds)
	})

	t.Run("includes user ids in the returned primitives", func(t *testing.T) {
		t.Parallel()
		svc, repo := newService(t)
		adminId := ids2.NewAdminId()
		userId := ids2.NewUserId()
		admin := anAdminWithUser(t, adminId, userId)

		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)

		result, err := svc.Get(context.Background(), get.NewQuery(adminId))
		require.NoError(t, err)

		prims, ok := result.(admins.AdminPrimitives)
		require.True(t, ok)
		assert.Equal(t, adminId.String(), prims.Id)
		assert.Contains(t, prims.UserIds, userId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *adminmocks.MockRepository, adminId ids2.AdminId)
	}{
		{
			desc: "propagates error when admin is not found in repo",
			setupRepo: func(repo *adminmocks.MockRepository, adminId ids2.AdminId) {
				repo.EXPECT().SearchById(mock.Anything, adminId).Return(nil, fmt.Errorf("db error"))
			},
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			svc, repo := newService(t)
			adminId := ids2.NewAdminId()
			tc.setupRepo(repo, adminId)

			_, err := svc.Get(context.Background(), get.NewQuery(adminId))
			assert.Error(t, err)
		})
	}
}
