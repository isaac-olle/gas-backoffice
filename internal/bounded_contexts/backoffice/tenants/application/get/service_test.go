package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns tenant when found and caller is self", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		userId := ids.NewUserId()
		tenant := aTenantWithId(t, userId)

		repo.EXPECT().SearchById(mock.Anything, userId).Return(tenant, nil)

		self := &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
		svc := get.NewService(domain_services.NewTenantFinder(repo))
		result, err := svc.Get(context.Background(), get.NewQuery(userId, self))
		require.NoError(t, err)

		ten, ok := result.(*tenants.Tenant)
		require.True(t, ok)
		assert.Equal(t, userId.String(), ten.ToPrimitives().Id)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		userId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, userId).Return(nil, fmt.Errorf("db error"))

		self := &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
		svc := get.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(userId, self))
		assert.Error(t, err)
	})

	t.Run("returns error when tenant not found", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		userId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, userId).Return(nil, nil)

		self := &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
		svc := get.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(userId, self))
		assert.Error(t, err)
	})

	t.Run("returns forbidden when tenant tries to access another tenant", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		userId := ids.NewUserId()
		tenant := aTenantWithId(t, userId)

		repo.EXPECT().SearchById(mock.Anything, userId).Return(tenant, nil)

		other := &authors.UserAuthor{Id: ids.NewUserId(), Role: vo.TenantRole}
		svc := get.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(userId, other))
		assert.Error(t, err)
	})

	t.Run("returns tenant when caller is the owning admin", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		userId := ids.NewUserId()
		adminId := ids.NewUserId()
		tenant := aTenantWithIdAndAdmin(t, userId, adminId)

		repo.EXPECT().SearchById(mock.Anything, userId).Return(tenant, nil)

		admin := &authors.UserAuthor{Id: adminId, Role: vo.AdminRole}
		svc := get.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(userId, admin))
		require.NoError(t, err)
	})
}
