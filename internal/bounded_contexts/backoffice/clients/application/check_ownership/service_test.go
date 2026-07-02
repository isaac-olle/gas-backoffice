package check_ownership_test

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	clientmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership primitives when owner matches", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		adminId := ids2.NewAdminId()
		tenantId := ids2.NewTenantId()
		clientId := ids2.NewClientId()
		c := aClientWithOwnership(t, adminId, tenantId)

		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		result, err := check_ownership.NewService(repo).CheckOwnership(
			context.Background(),
			check_ownership.NewQuery(clientId.String(), adminId, tenantId),
		)
		require.NoError(t, err)
		p, ok := result.(clients.OwnershipPrimitives)
		require.True(t, ok)
		assert.Equal(t, adminId.String(), p.AdminId)
		assert.Equal(t, tenantId.String(), p.TenantId)
	})

	errCases := []struct {
		desc      string
		query     func(t *testing.T) *check_ownership.Query
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id format",
			query:    func(t *testing.T) *check_ownership.Query { return aQuery(t, withClientId("bad-id")) },
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc:  "returns not found when repository errors",
			query: func(t *testing.T) *check_ownership.Query { return aQuery(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc:  "returns forbidden when admin does not match",
			query: func(t *testing.T) *check_ownership.Query { return aQuery(t, withAdminId(ids2.NewAdminId())) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c := aClientWithOwnership(t, ids2.NewAdminId(), ids2.NewTenantId())
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrInvalidOwnership,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			_, err := check_ownership.NewService(repo).CheckOwnership(context.Background(), tc.query(t))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
