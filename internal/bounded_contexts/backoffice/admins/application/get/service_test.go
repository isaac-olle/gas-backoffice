package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	adminmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns admin when found and caller is self", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		adminId := ids.NewUserId()
		admin := anAdminWithId(t, adminId)

		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)

		svc := get.NewService(domain_services.NewAdminFinder(repo))
		result, err := svc.Get(context.Background(), get.NewQuery(adminId, adminId))
		require.NoError(t, err)

		got, ok := result.(*admins.Admin)
		require.True(t, ok)
		assert.Equal(t, adminId.String(), got.ToPrimitives().Id)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		adminId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, adminId).Return(nil, fmt.Errorf("db error"))

		svc := get.NewService(domain_services.NewAdminFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(adminId, adminId))
		assert.Error(t, err)
	})

	t.Run("returns error when admin not found", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		adminId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, adminId).Return(nil, nil)

		svc := get.NewService(domain_services.NewAdminFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(adminId, adminId))
		assert.Error(t, err)
	})

	t.Run("returns forbidden when caller is not the admin itself", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		adminId := ids.NewUserId()
		admin := anAdminWithId(t, adminId)

		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)

		svc := get.NewService(domain_services.NewAdminFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(adminId, ids.NewUserId()))
		assert.Error(t, err)
	})
}
