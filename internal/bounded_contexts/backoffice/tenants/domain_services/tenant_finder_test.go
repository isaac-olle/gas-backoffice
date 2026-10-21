package domain_services_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestTenantFinder_Find(t *testing.T) {
	t.Parallel()

	t.Run("returns tenant when found", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		expected := aTenant(t)

		repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(expected, nil)

		finder := domain_services.NewTenantFinder(repo)
		got, err := finder.Find(context.Background(), ids2.NewUserId())

		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	errCases := []struct {
		desc       string
		findResult *tenants.Tenant
		findErr    error
		wantCode   app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrUserNotFound when tenant is nil",
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc:     "returns ErrUserNotFound when repository returns error",
			findErr:  fmt.Errorf("db error"),
			wantCode: tenants.ErrUserNotFound,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := tenantmocks.NewMockRepository(t)

			repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(tc.findResult, tc.findErr)

			finder := domain_services.NewTenantFinder(repo)
			_, err := finder.Find(context.Background(), ids2.NewUserId())

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

func aTenant(t *testing.T) *tenants.Tenant {
	t.Helper()
	userId := ids2.NewUserId()
	tenant, err := tenants.NewTenant(userId, "Test User", ids2.NewUserId(), "", false)
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
