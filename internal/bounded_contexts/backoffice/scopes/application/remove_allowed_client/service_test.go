package remove_allowed_client_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func newService(repo *scopemocks.MockRepository, bus *eventmocks.MockBus) *remove_allowed_client.Service {
	return remove_allowed_client.NewService(domain_services.NewScopeFinder(repo), repo, bus)
}

func TestRemoveAllowedClient(t *testing.T) {
	t.Parallel()

	t.Run("removes client from scope, saves updated scope, and publishes ScopeClientRemoved event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		prId := ids.NewProtectedResourceId()
		userId := ids.NewUserId()
		clientId := ids.NewClientId()
		scope, scopeId := aScopeWithClient(t, clientId, prId, userId)

		var saved *scopes.Scope
		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeClientRemovedType
			}))

		svc := newService(repo, bus)
		require.NoError(t, svc.RemoveAllowedClient(context.Background(),
			aCommand(t, withScopeId(scopeId.String()), withClientId(clientId.String()),
				withProtectedResourceId(prId.String()), withUserId(userId.String()))))

		assert.Empty(t, saved.ToPrimitives().AllowedClientIds)
	})

	prId := ids.NewProtectedResourceId()
	userId := ids.NewUserId()

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *remove_allowed_client.Command
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrScopeNotFound for invalid scope id",
			cmd:      func(t *testing.T) *remove_allowed_client.Command { return aCommand(t, withScopeId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc:     "returns ErrScopeNotFound for invalid client id",
			cmd:      func(t *testing.T) *remove_allowed_client.Command { return aCommand(t, withClientId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrScopeNotFound when scope does not exist",
			cmd:  func(t *testing.T) *remove_allowed_client.Command { return aCommand(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*scopes.Scope)(nil), fmt.Errorf("not found"))
			},
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrScopeNotSaved when save fails",
			cmd: func(t *testing.T) *remove_allowed_client.Command {
				return aCommand(t, withProtectedResourceId(prId.String()), withUserId(userId.String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				clientId := ids.NewClientId()
				scope, _ := aScopeWithClient(t, clientId, prId, userId)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: scopes.ErrScopeNotSaved,
		},
		{
			desc: "returns conflict when repository reports concurrent modification",
			cmd: func(t *testing.T) *remove_allowed_client.Command {
				return aCommand(t, withProtectedResourceId(prId.String()), withUserId(userId.String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				clientId := ids.NewClientId()
				scope, _ := aScopeWithClient(t, clientId, prId, userId)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).
					Return(persistance.ErrConcurrentModification)
			},
			wantCode: scopes.ErrScopeConcurrentModification,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := newService(repo, bus)
			err := svc.RemoveAllowedClient(context.Background(), tc.cmd(t))

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
