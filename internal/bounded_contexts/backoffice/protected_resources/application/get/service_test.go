package get_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*get.Service, *prmocks.MockRepository) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	finder := domain_services.NewProtectedResourceFinder(repo)
	return get.NewService(finder), repo
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns protected resource primitives with correct id, name and no scopes", func(t *testing.T) {
		t.Parallel()
		svc, repo := newService(t)
		prId := ids2.NewProtectedResourceId()
		pr := aPR(t, prId)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		result, err := svc.Get(context.Background(), get.NewQuery(prId))
		require.NoError(t, err)

		pr, ok := result.(*protected_resources.ProtectedResource)
		require.True(t, ok)
		prims := pr.ToPrimitives()
		assert.Equal(t, prId.String(), prims.Id)
		assert.Equal(t, "My Resource", prims.Name)
		assert.Empty(t, prims.ScopeIds)
		assert.NotEmpty(t, prims.TenantId)
		assert.NotEmpty(t, prims.AdminId)
	})

	t.Run("includes scope ids in the returned primitives", func(t *testing.T) {
		t.Parallel()
		svc, repo := newService(t)
		prId := ids2.NewProtectedResourceId()
		pr := aPR(t, prId)
		scopeId := ids2.NewScopeId()
		pr.AddScopeId(scopeId)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		result, err := svc.Get(context.Background(), get.NewQuery(prId))
		require.NoError(t, err)

		pr, ok := result.(*protected_resources.ProtectedResource)
		require.True(t, ok)
		prims := pr.ToPrimitives()
		assert.Equal(t, prId.String(), prims.Id)
		assert.Contains(t, prims.ScopeIds, scopeId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *prmocks.MockRepository) ids2.ProtectedResourceId
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns not found error when resource is not in repo",
			setupRepo: func(repo *prmocks.MockRepository) ids2.ProtectedResourceId {
				prId := ids2.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("not found"))
				return prId
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns not found error when repo returns nil",
			setupRepo: func(repo *prmocks.MockRepository) ids2.ProtectedResourceId {
				prId := ids2.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, nil)
				return prId
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := prmocks.NewMockRepository(t)
			finder := domain_services.NewProtectedResourceFinder(repo)
			prId := tc.setupRepo(repo)
			_, err := get.NewService(finder).Get(context.Background(), get.NewQuery(prId))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the get_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
