package add_allowed_client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func newService(repo *scopemocks.MockRepository, bus *eventmocks.MockBus, clientRepo *clientProjectionRepoStub) *add_allowed_client.Service {
	return add_allowed_client.NewService(domain_services.NewScopeFinder(repo), repo, bus, domain_services.NewClientProjectionFinder(clientRepo), sharedDomainServices.NewSubscriptionAuthorizer())
}

func TestAddAllowedClient(t *testing.T) {
	t.Parallel()

	t.Run("adds client to scope, saves updated scope, and publishes ScopeClientAdded event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		prId := ids.NewProtectedResourceId()
		subscriptionId := ids.NewSubscriptionId()
		scope, scopeId := aScope(t, prId, subscriptionId)
		clientId := ids.NewClientId()
		clientRepo := &clientProjectionRepoStub{projection: &scopes.ClientProjection{ClientId: clientId, SubscriptionId: subscriptionId}}

		var saved *scopes.Scope
		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeClientAddedType
			}))

		svc := newService(repo, bus, clientRepo)
		require.NoError(t, svc.AddAllowedClient(context.Background(),
			aCommand(t, withScopeId(scopeId.String()), withClientId(clientId.String()), withAuthor(authorFor(t, subscriptionId)))))

		assert.Contains(t, saved.ToPrimitives().AllowedClientIds, clientId.String())
	})

	prId := ids.NewProtectedResourceId()
	subscriptionId := ids.NewSubscriptionId()

	errCases := []struct {
		desc       string
		cmd        func(t *testing.T) *add_allowed_client.Command
		setupRepo  func(repo *scopemocks.MockRepository)
		clientRepo *clientProjectionRepoStub
		wantCode   app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrScopeNotFound for invalid scope id",
			cmd:      func(t *testing.T) *add_allowed_client.Command { return aCommand(t, withScopeId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrScopeNotFound when scope does not exist",
			cmd:  func(t *testing.T) *add_allowed_client.Command { return aCommand(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*scopes.Scope)(nil), fmt.Errorf("not found"))
			},
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns forbidden when author lacks subscription access",
			cmd: func(t *testing.T) *add_allowed_client.Command {
				return aCommand(t, withAuthor(authorFor(t, ids.NewSubscriptionId())))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t, prId, subscriptionId)
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(scope, nil)
			},
			wantCode: scopes.ErrInvalidOwnership,
		},
		{
			desc: "returns ErrScopeNotFound for invalid client id",
			cmd: func(t *testing.T) *add_allowed_client.Command {
				return aCommand(t, withClientId("not-a-uuid"), withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t, prId, subscriptionId)
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(scope, nil)
			},
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns not found when client projection lookup fails",
			cmd: func(t *testing.T) *add_allowed_client.Command {
				return aCommand(t, withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t, prId, subscriptionId)
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(scope, nil)
			},
			clientRepo: &clientProjectionRepoStub{err: fmt.Errorf("db error")},
			wantCode:   "CLIENT_PROJECTION_NOT_FOUND",
		},
		{
			desc: "returns forbidden when client and scope belong to different subscriptions",
			cmd: func(t *testing.T) *add_allowed_client.Command {
				return aCommand(t, withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t, prId, subscriptionId)
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(scope, nil)
			},
			clientRepo: &clientProjectionRepoStub{projection: &scopes.ClientProjection{ClientId: ids.NewClientId(), SubscriptionId: ids.NewSubscriptionId()}},
			wantCode:   scopes.ErrInvalidOwnership,
		},
		{
			desc: "returns ErrScopeNotSaved when save fails",
			cmd: func(t *testing.T) *add_allowed_client.Command {
				return aCommand(t, withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t, prId, subscriptionId)
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
			cmd: func(t *testing.T) *add_allowed_client.Command {
				return aCommand(t, withAuthor(authorFor(t, subscriptionId)))
			},
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t, prId, subscriptionId)
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
			clientRepo := tc.clientRepo
			if clientRepo == nil {
				clientRepo = &clientProjectionRepoStub{projection: &scopes.ClientProjection{ClientId: ids.NewClientId(), SubscriptionId: subscriptionId}}
			}

			svc := newService(repo, bus, clientRepo)
			err := svc.AddAllowedClient(context.Background(), tc.cmd(t))

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
