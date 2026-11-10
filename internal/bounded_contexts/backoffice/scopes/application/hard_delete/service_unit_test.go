package hard_delete_test

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func newService(finder *domain_services.ScopeFinder, repo *scopemocks.MockRepository, bus *eventmocks.MockBus) *hard_delete.Service {
	return hard_delete.NewService(finder, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the found scope with correct id and publishes ScopeDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		scope, scopeId, _, subscriptionId := aScope(t)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().HardDelete(mock.Anything, scopeId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeDeletedType
			}))

		finder := domain_services.NewScopeFinder(repo)
		svc := newService(finder, repo, bus)

		cmd := aCommand(t, withScopeId(scopeId), withAuthor(authorFor(t, subscriptionId)))
		require.NoError(t, svc.HardDelete(context.Background(), cmd))
		assert.Equal(t, scopeId.String(), scope.ToPrimitives().Id)
	})

	t.Run("returns error when finder returns error", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		_, scopeId, _, subscriptionId := aScope(t)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewScopeFinder(repo)
		svc := newService(finder, repo, bus)
		cmd := aCommand(t, withScopeId(scopeId), withAuthor(authorFor(t, subscriptionId)))
		err := svc.HardDelete(context.Background(), cmd)

		assert.Error(t, err)
	})

	t.Run("returns forbidden when author lacks subscription access", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		scope, scopeId, _, _ := aScope(t)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)

		finder := domain_services.NewScopeFinder(repo)
		svc := newService(finder, repo, bus)
		cmd := aCommand(t, withScopeId(scopeId), withAuthor(authorFor(t, ids2.NewSubscriptionId())))
		err := svc.HardDelete(context.Background(), cmd)

		assertErrCode(t, err, scopes.ErrInvalidOwnership)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *scopemocks.MockRepository) *hard_delete.Command
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns forbidden when scope is system resource",
			setupRepo: func(repo *scopemocks.MockRepository) *hard_delete.Command {
				scopeId := ids2.NewScopeId()
				prId := ids2.NewProtectedResourceId()
				subscriptionId := ids2.NewSubscriptionId()
				scope, err := scopes.FromPrimitives(scopes.ScopePrimitives{
					Id:                  scopeId.String(),
					Name:                "system-scope",
					ProtectedResourceId: prId.String(),
					SubscriptionId:      subscriptionId.String(),
					IsSystem:            true,
				})
				require.NoError(t, err)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
				return aCommand(t, withScopeId(scopeId), withAuthor(authorFor(t, subscriptionId)))
			},
			wantCode: scopes.ErrSystemResource,
		},
		{
			desc: "returns ErrScopeNotSaved when hard delete fails",
			setupRepo: func(repo *scopemocks.MockRepository) *hard_delete.Command {
				scope, scopeId, _, subscriptionId := aScope(t)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
				repo.EXPECT().
					HardDelete(mock.Anything, mock.Anything).
					Return(fmt.Errorf("db error"))
				return aCommand(t, withScopeId(scopeId), withAuthor(authorFor(t, subscriptionId)))
			},
			wantCode: scopes.ErrScopeNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			cmd := tc.setupRepo(repo)

			finder := domain_services.NewScopeFinder(repo)
			svc := newService(finder, repo, bus)
			err := svc.HardDelete(context.Background(), cmd)

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
