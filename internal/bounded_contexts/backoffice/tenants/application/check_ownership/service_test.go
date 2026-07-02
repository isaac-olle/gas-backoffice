package check_ownership_test

import (
	"context"
	app_errors2 "github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/check_ownership"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership primitives with correct admin when ownership matches", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		adminId := ids2.NewAdminId()
		tenantId := ids2.NewTenantId()
		tenant := aTenant(t, tenantId, adminId)

		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)

		svc := check_ownership.NewService(domain_services.NewTenantFinder(repo))
		result, err := svc.CheckOwnership(context.Background(), aQuery(t,
			withTenantId(tenantId.String()),
			withAdminId(adminId),
		))

		require.NoError(t, err)
		ownership, ok := result.(tenants.Ownership)
		require.True(t, ok)
		assert.Equal(t, adminId.String(), ownership.AdminId.String())
		assert.NotEmpty(t, ownership.AdminId.String())
	})

	errCases := []struct {
		desc      string
		qry       func(t *testing.T) *check_ownership.Query
		setupRepo func(repo *tenantmocks.MockRepository)
		wantType  app_errors2.ErrorType
	}{
		{
			desc:     "returns bad request for invalid tenant id format",
			qry:      func(t *testing.T) *check_ownership.Query { return aQuery(t, withTenantId("bad-tenant-id")) },
			wantType: app_errors2.BadRequestErrorType,
		},
		{
			desc: "returns not found when tenant does not exist",
			qry: func(t *testing.T) *check_ownership.Query {
				return aQuery(t, withTenantId(ids2.NewTenantId().String()), withAdminId(ids2.NewAdminId()))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), nil)
			},
			wantType: app_errors2.NotFoundErrorType,
		},
		{
			desc: "returns forbidden when admin does not own the tenant",
			qry: func(t *testing.T) *check_ownership.Query {
				tenantId := ids2.NewTenantId()
				return aQuery(t, withTenantId(tenantId.String()), withAdminId(ids2.NewAdminId()))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t, ids2.NewTenantId(), ids2.NewAdminId()), nil) // different admin
			},
			wantType: app_errors2.ForbiddenErrorType,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := tenantmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := check_ownership.NewService(domain_services.NewTenantFinder(repo))
			_, err := svc.CheckOwnership(context.Background(), tc.qry(t))

			assertErrType(t, err, tc.wantType)
		})
	}
}

// --- helpers ---

func assertErrType(t *testing.T, err error, errType app_errors2.ErrorType) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors2.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, errType, appErr.Type)
}
