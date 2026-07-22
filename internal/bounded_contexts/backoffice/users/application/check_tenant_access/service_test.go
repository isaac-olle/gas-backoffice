package check_tenant_access_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_tenant_access"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
	app_errors2 "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestCheckTenantAccess(t *testing.T) {
	t.Parallel()

	t.Run("admin who owns the tenant is granted access", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		tenantId := ids2.NewUserId()
		tenant := aTenant(t, tenantId, adminId)

		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)

		svc := check_tenant_access.NewService(domain_services.NewUserFinder(repo))
		result, err := svc.CheckTenantAccess(context.Background(), check_tenant_access.NewQuery(tenantId, anAdmin(adminId)))

		require.NoError(t, err)
		ownership, ok := result.(users.Ownership)
		require.True(t, ok)
		assert.Equal(t, tenantId.String(), ownership.UserId.String())
	})

	t.Run("admin who does not own the tenant is forbidden", func(t *testing.T) {
		t.Parallel()
		tenantId := ids2.NewUserId()
		tenant := aTenant(t, tenantId, ids2.NewUserId())

		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)

		svc := check_tenant_access.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.CheckTenantAccess(context.Background(), check_tenant_access.NewQuery(tenantId, anAdmin(ids2.NewUserId())))

		require.Error(t, err)
		var appErr *app_errors2.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, app_errors2.ForbiddenErrorType, appErr.Type)
	})

	t.Run("tenant accessing themselves is granted access", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		tenantId := ids2.NewUserId()
		tenant := aTenant(t, tenantId, adminId)

		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)

		svc := check_tenant_access.NewService(domain_services.NewUserFinder(repo))
		result, err := svc.CheckTenantAccess(context.Background(), check_tenant_access.NewQuery(tenantId, aTenantAuthor(tenantId)))

		require.NoError(t, err)
		ownership, ok := result.(users.Ownership)
		require.True(t, ok)
		assert.Equal(t, tenantId.String(), ownership.UserId.String())
	})

	t.Run("tenant accessing a different tenant is forbidden", func(t *testing.T) {
		t.Parallel()
		tenantId := ids2.NewUserId()
		tenant := aTenant(t, tenantId, ids2.NewUserId())

		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)

		svc := check_tenant_access.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.CheckTenantAccess(context.Background(), check_tenant_access.NewQuery(tenantId, aTenantAuthor(ids2.NewUserId())))

		require.Error(t, err)
		var appErr *app_errors2.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, app_errors2.ForbiddenErrorType, appErr.Type)
	})

	t.Run("returns not found when tenant does not exist", func(t *testing.T) {
		t.Parallel()
		tenantId := ids2.NewUserId()

		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return((*users.User)(nil), nil)

		svc := check_tenant_access.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.CheckTenantAccess(context.Background(), check_tenant_access.NewQuery(tenantId, anAdmin(ids2.NewUserId())))

		require.Error(t, err)
		var appErr *app_errors2.AppError
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, app_errors2.NotFoundErrorType, appErr.Type)
	})
}
