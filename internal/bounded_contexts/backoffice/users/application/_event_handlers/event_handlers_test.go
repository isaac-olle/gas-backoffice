package _event_handlers

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
)

// ---- AddClientOnClientCreated ----

func TestAddClientOnClientCreated(t *testing.T) {
	t.Parallel()

	t.Run("adds client id to user, saves updated user, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()
		adminId := ids2.NewAdminId()
		clientId := ids2.NewClientId()
		user := aUser(t, userId, adminId)

		var saved *users.User
		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *users.User) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // AddClientId does not record domain events; empty slice is expected
		}))

		handler := NewAddClientOnClientCreatedEventHandler(repo, bus)
		evt := backoffice.NewClientCreated(clientId.String(), "Client", adminId.String(), userId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.Contains(t, saved.ToPrimitives().ClientIds, clientId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrUserNotFound when repository find fails",
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), fmt.Errorf("not found"))
			},
			wantCode: users.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotSaved when save fails",
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aUser(t, ids2.NewUserId(), ids2.NewAdminId()), nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*users.User")).
					Return(fmt.Errorf("save error"))
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

			userId := ids2.NewUserId()
			adminId := ids2.NewAdminId()
			clientId := ids2.NewClientId()

			handler := NewAddClientOnClientCreatedEventHandler(repo, bus)
			evt := backoffice.NewClientCreated(clientId.String(), "Client", adminId.String(), userId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---- AddProtectedResourceOnProtectedResourceCreated ----

func TestAddProtectedResourceOnProtectedResourceCreated(t *testing.T) {
	t.Parallel()

	t.Run("adds protected resource id to user, saves updated user, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()
		adminId := ids2.NewAdminId()
		prId := ids2.NewProtectedResourceId()
		user := aUser(t, userId, adminId)

		var saved *users.User
		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *users.User) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // AddProtectedResourceId does not record domain events; empty slice is expected
		}))

		handler := NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(repo, bus)
		evt := backoffice.NewProtectedResourceCreated(prId.String(), "PR", userId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.Contains(t, saved.ToPrimitives().ProtectedResourceIds, prId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrUserNotFound when repository find fails",
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), fmt.Errorf("not found"))
			},
			wantCode: users.ErrUserNotFound,
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

			userId := ids2.NewUserId()
			prId := ids2.NewProtectedResourceId()

			handler := NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(repo, bus)
			evt := backoffice.NewProtectedResourceCreated(prId.String(), "PR", userId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---- RemoveClientOnClientDeleted ----

func TestRemoveClientOnClientDeleted(t *testing.T) {
	t.Parallel()

	t.Run("removes client id from user, saves updated user, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()
		adminId := ids2.NewAdminId()
		clientId := ids2.NewClientId()

		prims := users.UserPrimitives{
			Id:        userId.String(),
			Name:      "Test User",
			AdminId:   adminId.String(),
			ClientIds: []string{clientId.String()},
		}
		user, err := users.FromPrimitives(prims)
		require.NoError(t, err)

		var saved *users.User
		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *users.User) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // RemoveClientId does not record domain events; empty slice is expected
		}))

		handler := NewRemoveClientOnClientDeletedEventHandler(repo, bus)
		evt := backoffice.NewClientDeleted(clientId.String(), userId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))
		assert.Empty(t, saved.ToPrimitives().ClientIds)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrUserNotFound when repository find fails",
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), fmt.Errorf("db error"))
			},
			wantCode: users.ErrUserNotFound,
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

			userId := ids2.NewUserId()
			clientId := ids2.NewClientId()

			handler := NewRemoveClientOnClientDeletedEventHandler(repo, bus)
			evt := backoffice.NewClientDeleted(clientId.String(), userId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---- RemoveProtectedResourceOnProtectedResourceDeleted ----

func TestRemoveProtectedResourceOnProtectedResourceDeleted(t *testing.T) {
	t.Parallel()

	t.Run("removes protected resource id from user, saves updated user, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()
		adminId := ids2.NewAdminId()
		prId := ids2.NewProtectedResourceId()

		prims := users.UserPrimitives{
			Id:                   userId.String(),
			Name:                 "Test User",
			AdminId:              adminId.String(),
			ProtectedResourceIds: []string{prId.String()},
		}
		user, err := users.FromPrimitives(prims)
		require.NoError(t, err)

		var saved *users.User
		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *users.User) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // RemoveProtectedResourceId does not record domain events; empty slice is expected
		}))

		handler := NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backoffice.NewProtectedResourceDeleted(prId.String(), userId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))
		assert.Empty(t, saved.ToPrimitives().ProtectedResourceIds)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *usermocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrUserNotFound when repository find fails",
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*users.User)(nil), fmt.Errorf("db error"))
			},
			wantCode: users.ErrUserNotFound,
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

			userId := ids2.NewUserId()
			prId := ids2.NewProtectedResourceId()

			handler := NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(repo, bus)
			evt := backoffice.NewProtectedResourceDeleted(prId.String(), userId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

func aUser(t *testing.T, userId ids2.UserId, adminId ids2.AdminId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(adminId))
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
