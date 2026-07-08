package hard_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
)

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the found user and publishes UserDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		userId := ids.NewUserId()
		user := aUserWithId(t, userId)
		cmd := aCommand(t, withUserId(userId.String()))

		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)
		repo.EXPECT().HardDelete(mock.Anything, userId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserDeletedType
			}))

		svc := hard_delete.NewService(domain_services.NewUserFinder(repo), repo, bus)
		require.NoError(t, svc.HardDelete(context.Background(), cmd))
		assert.Equal(t, userId.String(), user.ToPrimitives().Id)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *hard_delete.Command
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrUserNotFound for invalid user id format",
			cmd:      func(t *testing.T) *hard_delete.Command { return aCommand(t, withUserId("bad-id")) },
			wantCode: users.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when user does not exist",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), nil)
			},
			wantCode: users.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when repository find returns error",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), fmt.Errorf("db error"))
			},
			wantCode: users.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotSaved when hard delete fails",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aUser(t), nil)
				repo.EXPECT().
					HardDelete(mock.Anything, mock.Anything).
					Return(fmt.Errorf("db error"))
			},
			wantCode: users.ErrUserNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := usermocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := hard_delete.NewService(domain_services.NewUserFinder(repo), repo, bus)
			err := svc.HardDelete(context.Background(), tc.cmd(t))

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
