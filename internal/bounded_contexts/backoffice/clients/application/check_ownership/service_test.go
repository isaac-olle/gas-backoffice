package check_ownership_test

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership primitives when owner matches", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		userId := ids2.NewUserId()
		clientId := ids2.NewClientId()
		c := aClientWithOwnership(t, userId)

		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)

		result, err := check_ownership.NewService(repo).CheckOwnership(
			context.Background(),
			check_ownership.NewQuery(clientId.String(), userId.String()),
		)
		require.NoError(t, err)
		p, ok := result.(clients.OwnershipPrimitives)
		require.True(t, ok)
		assert.Equal(t, userId.String(), p.UserId)
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
			desc:  "returns forbidden when user does not match",
			query: func(t *testing.T) *check_ownership.Query { return aQuery(t, withUserId(ids2.NewUserId())) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c := aClientWithOwnership(t, ids2.NewUserId())
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
