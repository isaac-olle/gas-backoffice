package check_ownership_test

import (
	"context"
	"fmt"
	app_errors2 "github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/check_ownership"
	scopemocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership primitives with all correct ids when ownership matches", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		scopeId := ids2.NewScopeId()
		prId := ids2.NewProtectedResourceId()
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		scope := aScope(t, scopeId, prId, tenantId, adminId)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)

		svc := check_ownership.NewService(repo)
		result, err := svc.CheckOwnership(context.Background(), aQuery(t,
			withScopeId(scopeId.String()),
			withProtectedResourceId(prId.String()),
			withTenantId(tenantId.String()),
			withAdminId(adminId.String()),
		))

		require.NoError(t, err)
		ownership, ok := result.(scopes.Ownership)
		require.True(t, ok)
		prims := ownership.ToPrimitives()
		assert.Equal(t, prId.String(), prims.ProtectedResourceId)
		assert.Equal(t, tenantId.String(), prims.TenantId)
		assert.Equal(t, adminId.String(), prims.AdminId)
	})

	errCases := []struct {
		desc      string
		qry       func(t *testing.T) *check_ownership.Query
		setupRepo func(repo *scopemocks.MockRepository)
		wantType  app_errors2.ErrorType
	}{
		{
			desc:     "returns bad request for invalid scope id",
			qry:      func(t *testing.T) *check_ownership.Query { return aQuery(t, withScopeId("bad-id")) },
			wantType: app_errors2.BadRequestErrorType,
		},
		{
			desc:     "returns bad request when protected resource id is invalid",
			qry:      func(t *testing.T) *check_ownership.Query { return aQuery(t, withProtectedResourceId("bad-id")) },
			wantType: app_errors2.BadRequestErrorType,
		},
		{
			desc: "returns not found when scope does not exist",
			qry:  func(t *testing.T) *check_ownership.Query { return aQuery(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*scopes.Scope)(nil), fmt.Errorf("not found"))
			},
			wantType: app_errors2.NotFoundErrorType,
		},
		{
			desc: "returns forbidden when protected resource id does not match",
			qry: func(t *testing.T) *check_ownership.Query {
				// query uses a different protected resource id so ownership check fails
				return aQuery(t, withProtectedResourceId(ids2.NewProtectedResourceId().String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scopeId := ids2.NewScopeId()
				prId := ids2.NewProtectedResourceId()
				tenantId := ids2.NewTenantId()
				adminId := ids2.NewAdminId()
				scope := aScope(t, scopeId, prId, tenantId, adminId)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
			},
			wantType: app_errors2.ForbiddenErrorType,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := check_ownership.NewService(repo)
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
