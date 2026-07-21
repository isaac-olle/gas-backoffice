package check_ownership_test

import (
	"context"
	app_errors2 "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
)

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns ownership when ownership matches", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		userId := ids2.NewUserId()
		user := aUser(t, userId)

		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)

		svc := check_ownership.NewService(domain_services.NewUserFinder(repo))
		result, err := svc.CheckOwnership(context.Background(), aQuery(t,
			withUserId(userId.String()),
		))

		require.NoError(t, err)
		ownership, ok := result.(users.Ownership)
		require.True(t, ok)
		assert.Equal(t, userId.String(), ownership.UserId.String())
	})

	errCases := []struct {
		desc      string
		qry       func(t *testing.T) *check_ownership.Query
		setupRepo func(repo *usermocks.MockRepository)
		wantType  app_errors2.ErrorType
	}{
		{
			desc:     "returns bad request for invalid user id format",
			qry:      func(t *testing.T) *check_ownership.Query { return aQuery(t, withUserId("bad-user-id")) },
			wantType: app_errors2.BadRequestErrorType,
		},
		{
			desc: "returns not found when user does not exist",
			qry: func(t *testing.T) *check_ownership.Query {
				return aQuery(t, withUserId(ids2.NewUserId().String()))
			},
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), nil)
			},
			wantType: app_errors2.NotFoundErrorType,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := usermocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := check_ownership.NewService(domain_services.NewUserFinder(repo))
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
