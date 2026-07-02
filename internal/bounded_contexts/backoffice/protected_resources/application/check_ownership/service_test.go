package check_ownership_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	prmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*check_ownership.Service, *prmocks.MockRepository) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	return check_ownership.NewService(repo), repo
}

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership with matching tenant and admin ids", func(t *testing.T) {
		t.Parallel()
		svc, repo := newService(t)
		prId := ids2.NewProtectedResourceId()
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		pr := aPRWithOwnership(t, prId, tenantId, adminId)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		result, err := svc.CheckOwnership(context.Background(), check_ownership.NewQuery(prId.String(), tenantId.String(), adminId))
		require.NoError(t, err)
		require.NotNil(t, result)

		ownership, ok := result.(protected_resources.Ownership)
		require.True(t, ok)
		prims := ownership.ToPrimitives()
		assert.Equal(t, tenantId.String(), prims.TenantId)
		assert.Equal(t, adminId.String(), prims.AdminId)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *prmocks.MockRepository) (prIdStr string, tenantIdStr string, adminId ids2.AdminId)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns bad request error for invalid tenant id",
			setupRepo: func(_ *prmocks.MockRepository) (string, string, ids2.AdminId) {
				return ids2.NewProtectedResourceId().String(), "not-a-uuid", ids2.NewAdminId()
			},
			wantCode: protected_resources.ErrInvalidTenantId,
		},
		{
			desc: "returns bad request error for invalid protected resource id",
			setupRepo: func(_ *prmocks.MockRepository) (string, string, ids2.AdminId) {
				return "bad-id", ids2.NewTenantId().String(), ids2.NewAdminId()
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns not found error when resource is not in repo",
			setupRepo: func(repo *prmocks.MockRepository) (string, string, ids2.AdminId) {
				prId := ids2.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("not found"))
				return prId.String(), ids2.NewTenantId().String(), ids2.NewAdminId()
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns forbidden error when ownership does not match",
			setupRepo: func(repo *prmocks.MockRepository) (string, string, ids2.AdminId) {
				prId := ids2.NewProtectedResourceId()
				tenantId := ids2.NewTenantId()
				adminId := ids2.NewAdminId()
				pr := aPRWithOwnership(t, prId, tenantId, adminId)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return prId.String(), ids2.NewTenantId().String(), adminId // different tenant
			},
			wantCode: protected_resources.ErrInvalidOwnership,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			svc, repo := newService(t)
			prIdStr, tenantIdStr, adminId := tc.setupRepo(repo)
			_, err := svc.CheckOwnership(context.Background(), check_ownership.NewQuery(prIdStr, tenantIdStr, adminId))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the check_ownership_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
