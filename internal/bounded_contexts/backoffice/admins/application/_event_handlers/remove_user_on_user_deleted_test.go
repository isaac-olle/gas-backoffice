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

// makeAdminWithUser builds an Admin that already contains the given user.
func makeAdminWithUser(t *testing.T, adminId ids2.AdminId, userId ids2.UserId) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:      adminId.String(),
		Name:    "Test",
		UserIds: []string{userId.String()},
	})
	require.NoError(t, err)
	return a
}

func newRemoveHandler(repo admins.Repository, bus *eventmocks.MockBus) *_event_handlers.RemoveUserOnUserDeletedEventHandler {
	finder := domain_services.NewAdminFinder(repo)
	return _event_handlers.NewRemoveUserOnUserDeletedEventHandler(repo, finder, bus)
}

func TestRemoveUserOnUserDeleted(t *testing.T) {
	t.Parallel()

	t.Run("saves admin with user id removed and publishes empty events", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		userId := ids2.NewUserId()
		admin := makeAdminWithUser(t, adminId, userId)

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

		evt := backofficeevents.NewUserDeleted(userId.String(), adminId.String())
		require.NoError(t, newRemoveHandler(repo, bus).Handle(context.Background(), evt))

		require.NotNil(t, savedAdmin)
		prims := savedAdmin.ToPrimitives()
		assert.Equal(t, adminId.String(), prims.Id)
		assert.NotContains(t, prims.UserIds, userId.String())
	})

	t.Run("still saves when user id is not present in admin", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		userId := ids2.NewUserId()
		otherUserId := ids2.NewUserId()
		admin := makeAdminWithUser(t, adminId, otherUserId)

		repo := adminmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().SearchById(mock.Anything, adminId).Return(admin, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(a *admins.Admin) bool {
				return !containsStr(a.ToPrimitives().UserIds, userId.String())
			})).
			Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return len(evts) == 0
		})).Return()

		evt := backofficeevents.NewUserDeleted(userId.String(), adminId.String())
		require.NoError(t, newRemoveHandler(repo, bus).Handle(context.Background(), evt))
	})

	userId := ids2.NewUserId()
	errCases := []struct {
		desc      string
		adminId   string
		setupRepo func(repo *adminmocks.MockRepository, adminId ids2.AdminId, userId ids2.UserId)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns bad request error when admin id in event is invalid",
			adminId:  "not-a-uuid",
			wantCode: admins.ErrAdminNotFound,
		},
		{
			desc: "returns error when admin is not found in repo",
			setupRepo: func(repo *adminmocks.MockRepository, adminId ids2.AdminId, _ ids2.UserId) {
				repo.EXPECT().SearchById(mock.Anything, adminId).Return(nil, fmt.Errorf("db error"))
			},
		},
		{
			desc: "returns internal error when save fails",
			setupRepo: func(repo *adminmocks.MockRepository, adminId ids2.AdminId, userId ids2.UserId) {
				admin := makeAdminWithUser(t, adminId, userId)
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
				tc.setupRepo(repo, adminId, userId)
			}
			evt := backofficeevents.NewUserDeleted(userId.String(), adminIdStr)
			err := newRemoveHandler(repo, bus).Handle(context.Background(), evt)
			if tc.wantCode != "" {
				assertErrCode(t, err, tc.wantCode)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// containsStr reports whether s is in slice.
func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
