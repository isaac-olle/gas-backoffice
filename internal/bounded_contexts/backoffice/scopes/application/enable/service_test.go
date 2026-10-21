package enable_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/enable"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func newService(repo *scopemocks.MockRepository, bus *eventmocks.MockBus) *enable.Service {
	return enable.NewService(domain_services.NewScopeFinder(repo), repo, bus)
}

func TestEnable(t *testing.T) {
	t.Parallel()

	t.Run("enables the scope, saves it, and publishes ScopeEnabled event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		scopeId := ids2.NewScopeId()
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()
		scope := aDisabledScope(t, scopeId, prId, userId)

		var saved *scopes.Scope
		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeEnabledType
			}))

		svc := newService(repo, bus)
		require.NoError(t, svc.Enable(context.Background(),
			aCommand(t, withScopeId(scopeId.String()), withProtectedResourceId(prId.String()), withUserId(userId.String()))))

		assert.True(t, saved.IsEnabled())
	})

	scopeId := ids2.NewScopeId()
	prId := ids2.NewProtectedResourceId()
	userId := ids2.NewUserId()

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *enable.Command
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrScopeNotFound for invalid scope id format",
			cmd:      func(t *testing.T) *enable.Command { return aCommand(t, withScopeId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrInvalidOwnership when ownership does not match",
			cmd: func(t *testing.T) *enable.Command {
				return aCommand(t, withScopeId(scopeId.String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope := aDisabledScope(t, scopeId, prId, userId)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
			},
			wantCode: scopes.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when scope is a system resource",
			cmd: func(t *testing.T) *enable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withProtectedResourceId(prId.String()), withUserId(userId.String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, err := scopes.FromPrimitives(scopes.ScopePrimitives{
					Id:                  scopeId.String(),
					Name:                "system-scope",
					ProtectedResourceId: prId.String(),
					AdminId:             ids2.NewUserId().String(),
					UserId:              userId.String(),
					IsSystem:            true,
				})
				require.NoError(t, err)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
			},
			wantCode: scopes.ErrSystemResource,
		},
		{
			desc: "returns ErrScopeNotSaved when save fails",
			cmd: func(t *testing.T) *enable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withProtectedResourceId(prId.String()), withUserId(userId.String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope := aDisabledScope(t, scopeId, prId, userId)
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
			cmd: func(t *testing.T) *enable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withProtectedResourceId(prId.String()), withUserId(userId.String()))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope := aDisabledScope(t, scopeId, prId, userId)
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
			err := svc.Enable(context.Background(), tc.cmd(t))

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
