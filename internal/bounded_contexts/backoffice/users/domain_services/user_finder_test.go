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

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
)

func TestUserFinder_Find(t *testing.T) {
	t.Parallel()

	t.Run("returns user when found", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		expected := aUser(t)

		repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(expected, nil)

		finder := domain_services.NewUserFinder(repo)
		got, err := finder.Find(context.Background(), ids2.NewUserId())

		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	errCases := []struct {
		desc       string
		findResult *users.User
		findErr    error
		wantCode   app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrUserNotFound when user is nil",
			wantCode: users.ErrUserNotFound,
		},
		{
			desc:     "returns ErrUserNotFound when repository returns error",
			findErr:  fmt.Errorf("db error"),
			wantCode: users.ErrUserNotFound,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := usermocks.NewMockRepository(t)

			repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(tc.findResult, tc.findErr)

			finder := domain_services.NewUserFinder(repo)
			_, err := finder.Find(context.Background(), ids2.NewUserId())

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

func aUser(t *testing.T) *users.User {
	t.Helper()
	userId := ids2.NewUserId()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(userId))
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
