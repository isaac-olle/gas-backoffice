package _event_handlers_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	_event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/_event_handlers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	adminmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Object Mother
// ---------------------------------------------------------------------------

// makeAdmin builds a minimal valid Admin with no users.
func makeAdmin(t *testing.T, adminId ids2.AdminId) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:      adminId.String(),
		Name:    "Test",
		UserIds: []string{},
	})
	require.NoError(t, err)
	return a
}

func newAddHandler(repo admins.Repository, bus *eventmocks.MockBus) *_event_handlers.AddUserOnUserCreatedEventHandler {
	finder := domain_services.NewAdminFinder(repo)
	return _event_handlers.NewAddUserOnUserCreatedEventHandler(repo, finder, bus)
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the whole _event_handlers_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestAddUserOnUserCreated(t *testing.T) {
	t.Parallel()

	t.Run("saves admin with user id added and publishes empty events", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		userId := ids2.NewUserId()
		admin := makeAdmin(t, adminId)

		var savedAdmin *admins.Admin
		repo := adminmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(a *admins.Admin) bool {
				savedAdmin = a
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 0
			})).
			Return()

		evt := backofficeevents.NewUserCreated(userId.String(), "My User", adminId.String())
		require.NoError(t, newAddHandler(repo, bus).Handle(context.Background(), evt))

		require.NotNil(t, savedAdmin)
		prims := savedAdmin.ToPrimitives()
		assert.Equal(t, adminId.String(), prims.Id)
		assert.Contains(t, prims.UserIds, userId.String())
	})

	userId := ids2.NewUserId()
	errCases := []struct {
		desc      string
		adminId   string
		setupRepo func(repo *adminmocks.MockRepository, adminId ids2.AdminId)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns bad request error when admin id in event is invalid",
			adminId:  "not-a-uuid",
			wantCode: admins.ErrAdminNotFound,
		},
		{
			desc: "returns error when admin is not found in repo",
			setupRepo: func(repo *adminmocks.MockRepository, adminId ids2.AdminId) {
				repo.EXPECT().SearchById(mock.Anything, adminId).Return(nil, fmt.Errorf("db error"))
			},
		},
		{
			desc: "returns internal error when save fails",
			setupRepo: func(repo *adminmocks.MockRepository, adminId ids2.AdminId) {
				admin := makeAdmin(t, adminId)
				repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*admins.Admin")).Return(fmt.Errorf("save error"))
			},
			wantCode: admins.ErrAdminNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := adminmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			adminId := ids2.NewAdminId()
			adminIdStr := adminId.String()
			if tc.adminId != "" {
				adminIdStr = tc.adminId
			}
			if tc.setupRepo != nil {
				tc.setupRepo(repo, adminId)
			}
			evt := backofficeevents.NewUserCreated(userId.String(), "My User", adminIdStr)
			err := newAddHandler(repo, bus).Handle(context.Background(), evt)
			if tc.wantCode != "" {
				assertErrCode(t, err, tc.wantCode)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
