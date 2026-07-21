package change_name_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
)

func TestChangeName(t *testing.T) {
	t.Parallel()

	t.Run("saves user with updated name and publishes UserNameChanged event", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		userId := ids2.NewUserId()
		user := aUser(t, userId, nil, "Old Name")

		var saved *users.User
		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *users.User) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserNameChangedType
			}))

		svc := change_name.NewService(domain_services.NewUserFinder(repo), repo, bus)
		require.NoError(t, svc.ChangeName(context.Background(), aCommand(t, withUserId(userId), withName("New Name"))))

		assert.Equal(t, "New Name", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *change_name.Command
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrUserNameNotChanged for empty name",
			cmd:      func(t *testing.T) *change_name.Command { return aCommand(t, withName("")) },
			wantCode: users.ErrUserNameNotChanged,
		},
		{
			desc:     "returns ErrUserNameNotChanged for reserved name",
			cmd:      func(t *testing.T) *change_name.Command { return aCommand(t, withName("Gas-User")) },
			wantCode: users.ErrUserNameNotChanged,
		},
		{
			desc: "returns ErrUserNotFound when user does not exist",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), nil)
			},
			wantCode: users.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when repository find returns error",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), fmt.Errorf("db error"))
			},
			wantCode: users.ErrUserNotFound,
		},
		{
			desc: "returns internal error when save fails",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aUser(t, ids2.NewUserId(), nil, "Old Name"), nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*users.User")).
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

			svc := change_name.NewService(domain_services.NewUserFinder(repo), repo, bus)
			err := svc.ChangeName(context.Background(), tc.cmd(t))

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
