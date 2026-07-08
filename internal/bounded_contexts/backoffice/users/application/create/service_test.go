package create_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves user with correct data and publishes UserCreated event", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		cmd := aCommand(t, withName("my-user"))

		var saved *users.User
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *users.User) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserCreatedType
			}))

		require.NoError(t, create.NewService(repo, bus).Create(context.Background(), cmd))
		assert.Equal(t, "my-user", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *create.Command
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects empty name with ErrUserNotCreated",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withName("")) },
			wantCode: users.ErrUserNotCreated,
		},
		{
			desc:     "rejects reserved name with ErrUserNotCreated",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withName("Gas-User")) },
			wantCode: users.ErrUserNotCreated,
		},
		{
			desc: "wraps repository error as ErrUserNotSaved",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *usermocks.MockRepository) {
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

			svc := create.NewService(repo, bus)
			err := svc.Create(context.Background(), tc.cmd(t))

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
