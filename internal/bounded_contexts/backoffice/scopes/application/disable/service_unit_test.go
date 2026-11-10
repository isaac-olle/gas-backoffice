package disable_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/disable"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func newService(repo *scopemocks.MockRepository, bus *eventmocks.MockBus) *disable.Service {
	return disable.NewService(domain_services.NewScopeFinder(repo), repo, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func TestDisable(t *testing.T) {
	t.Parallel()

	t.Run("disables the scope, saves it, and publishes ScopeDisabled event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		scopeId := ids2.NewScopeId()
		prId := ids2.NewProtectedResourceId()
		subscriptionId := ids2.NewSubscriptionId()
		scope := anEnabledScope(t, scopeId, prId, subscriptionId)

		var saved *scopes.Scope
		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeDisabledType
			}))

		svc := newService(repo, bus)
		require.NoError(t, svc.Disable(context.Background(),
			aCommand(t, withScopeId(scopeId.String()), withAuthor(authorFor(t, subscriptionId)))))

		assert.False(t, saved.IsEnabled())
	})

	scopeId := ids2.NewScopeId()
	prId := ids2.NewProtectedResourceId()
	subscriptionId := ids2.NewSubscriptionId()

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *disable.Command
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrScopeNotFound for invalid scope id format",
			cmd:      func(t *testing.T) *disable.Command { return aCommand(t, withScopeId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrInvalidOwnership when author lacks subscription access",
			cmd: func(t *testing.T) *disable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withAuthor(authorFor(t, ids2.NewSubscriptionId())))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope := anEnabledScope(t, scopeId, prId, subscriptionId)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
			},
			wantCode: scopes.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when scope is a system resource",
			cmd: func(t *testing.T) *disable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, err := scopes.FromPrimitives(scopes.ScopePrimitives{
					Id:                  scopeId.String(),
					Name:                "system-scope",
					ProtectedResourceId: prId.String(),
					SubscriptionId:      subscriptionId.String(),
					IsSystem:            true,
					Enabled:             true,
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
			cmd: func(t *testing.T) *disable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope := anEnabledScope(t, scopeId, prId, subscriptionId)
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
			cmd: func(t *testing.T) *disable.Command {
				return aCommand(t, withScopeId(scopeId.String()), withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope := anEnabledScope(t, scopeId, prId, subscriptionId)
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
			err := svc.Disable(context.Background(), tc.cmd(t))

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
