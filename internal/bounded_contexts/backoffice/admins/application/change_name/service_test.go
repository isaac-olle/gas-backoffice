package change_name_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/change_name"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	adminmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/mocks"
)

func TestChangeName(t *testing.T) {
	t.Parallel()

	t.Run("saves admin with updated name and publishes AdminNameChanged event", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		adminId := ids2.NewUserId()
		admin := anAdmin(t, adminId, "Old Name")

		var saved *admins.Admin
		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(a *admins.Admin) bool {
				saved = a
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.AdminNameChangedType
			}))

		svc := change_name.NewService(domain_services.NewAdminFinder(repo), repo, bus)
		require.NoError(t, svc.ChangeName(context.Background(), aCommand(t, withAdminId(adminId), withName("New Name"))))

		assert.Equal(t, "New Name", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *change_name.Command
		setupRepo func(repo *adminmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrAdminNameNotChanged for empty name",
			cmd:      func(t *testing.T) *change_name.Command { return aCommand(t, withName("")) },
			wantCode: admins.ErrAdminNameNotChanged,
		},
		{
			desc: "returns ErrAdminNotFound when admin does not exist",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *adminmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*admins.Admin)(nil), nil)
			},
			wantCode: admins.ErrAdminNotFound,
		},
		{
			desc: "returns ErrAdminNotFound when repository find returns error",
			cmd:  func(t *testing.T) *change_name.Command { return aCommand(t, withName("Valid Name")) },
			setupRepo: func(repo *adminmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*admins.Admin)(nil), fmt.Errorf("db error"))
			},
			wantCode: admins.ErrAdminNotFound,
		},
		{
			desc: "returns internal error when save fails",
			cmd: func(t *testing.T) *change_name.Command {
				adminId := ids2.NewUserId()
				return aCommand(t, withAdminId(adminId), withName("Valid Name"))
			},
			setupRepo: func(repo *adminmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids2.UserId) (*admins.Admin, error) {
						return anAdmin(t, id, "Old Name"), nil
					})
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*admins.Admin")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: admins.ErrAdminNotSaved,
		},
		{
			desc: "returns conflict when repository reports concurrent modification",
			cmd: func(t *testing.T) *change_name.Command {
				adminId := ids2.NewUserId()
				return aCommand(t, withAdminId(adminId), withName("Valid Name"))
			},
			setupRepo: func(repo *adminmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids2.UserId) (*admins.Admin, error) {
						return anAdmin(t, id, "Old Name"), nil
					})
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*admins.Admin")).
					Return(persistance.ErrConcurrentModification)
			},
			wantCode: admins.ErrAdminConcurrentModification,
		},
		{
			desc: "returns forbidden when caller is not the admin itself",
			cmd: func(t *testing.T) *change_name.Command {
				return aCommand(t, withName("Valid Name"), withCallerId(ids2.NewUserId()))
			},
			setupRepo: func(repo *adminmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids2.UserId) (*admins.Admin, error) {
						return anAdmin(t, id, "Old Name"), nil
					})
			},
			wantCode: admins.ErrInvalidOwnership,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := adminmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := change_name.NewService(domain_services.NewAdminFinder(repo), repo, bus)
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
