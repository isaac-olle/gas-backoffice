package domain_services_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// aPR builds a minimal valid ProtectedResource for test setup.
func aPR(t *testing.T) *protected_resources.ProtectedResource {
	t.Helper()
	ownership := protected_resources.NewOwnership(ids2.NewUserId(), ids2.NewUserId())
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "Test", ownership, false)
	require.NoError(t, err)
	pr.PullDomainEvents() // drain creation event
	return pr
}

func TestProtectedResourceFinder(t *testing.T) {
	t.Parallel()

	t.Run("returns protected resource when repo finds it", func(t *testing.T) {
		t.Parallel()
		prId := ids2.NewProtectedResourceId()
		expected := aPR(t)

		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, prId).Return(expected, nil)

		finder := domain_services.NewProtectedResourceFinder(repo)
		got, err := finder.Find(context.Background(), prId)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	errCases := []struct {
		desc     string
		repoRet  *protected_resources.ProtectedResource
		repoErr  error
		wantCode app_errors.ErrorCode
	}{
		{
			desc:     "returns internal error when repo returns an error",
			repoErr:  fmt.Errorf("db error"),
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc:     "returns not found error when repo returns nil",
			repoRet:  nil,
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			prId := ids2.NewProtectedResourceId()

			repo := prmocks.NewMockRepository(t)
			repo.EXPECT().SearchById(mock.Anything, prId).Return(tc.repoRet, tc.repoErr)

			finder := domain_services.NewProtectedResourceFinder(repo)
			_, err := finder.Find(context.Background(), prId)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the domain_services_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
