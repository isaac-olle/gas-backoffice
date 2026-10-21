package remove_allowed_grant_type_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_grant_type"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

func TestRemoveAllowedGrantType(t *testing.T) {
	t.Parallel()

	t.Run("removes grant type from client and publishes ClientGrantTypeRemoved event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		userId := ids2.NewUserId()
		c, clientId := aClientWithGrantType(t, userId, "client_credentials")

		var saved *clients.Client
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(cl *clients.Client) bool {
				saved = cl
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientGrantTypeRemovedType
			}))

		cmd := remove_allowed_grant_type.NewCommand(clientId, ids2.NewUserId(), userId, "client_credentials")
		require.NoError(t, newService(t, repo, bus).Remove(context.Background(), cmd))

		assert.NotContains(t, saved.ToPrimitives().AllowedGrantTypes, "client_credentials")
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *remove_allowed_grant_type.Command
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid grant type name",
			cmd:      func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t, withGrantType("implicit")) },
			wantCode: vos.ErrInvalidGrantType,
		},
		{
			desc: "returns not found when client does not exist",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "propagates repository search error",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns forbidden when ownership does not match",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClientWithGrantType(t, ids2.NewUserId(), "client_credentials")
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when client is system resource",
			cmd: func(t *testing.T) *remove_allowed_grant_type.Command {
				return aCommand(t, withUserId(systemResourceUserId))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.FromPrimitives(clients.ClientPrimitives{
					Id:                ids2.NewClientId().String(),
					Name:              "system-client",
					AdminId:           ids2.NewUserId().String(),
					UserId:            systemResourceUserId.String(),
					AllowedGrantTypes: []string{"client_credentials"},
					IsSystem:          true,
				})
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrSystemResource,
		},
		{
			desc: "wraps repository save failure",
			cmd: func(t *testing.T) *remove_allowed_grant_type.Command {
				return aCommand(t, withUserId(saveFailureUserId))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewUserId(), saveFailureUserId), nil, []string{"client_credentials"}, false, false)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
		{
			desc: "returns conflict when repository reports concurrent modification",
			cmd: func(t *testing.T) *remove_allowed_grant_type.Command {
				return aCommand(t, withUserId(concurrentModificationUserId))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewUserId(), concurrentModificationUserId), nil, []string{"client_credentials"}, false, false)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(persistance.ErrConcurrentModification)
			},
			wantCode: clients.ErrClientConcurrentModification,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo, bus := newMocks(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			assertErrCode(t, newService(t, repo, bus).Remove(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

var (
	systemResourceUserId         = ids2.NewUserId()
	saveFailureUserId            = ids2.NewUserId()
	concurrentModificationUserId = ids2.NewUserId()
)

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func newService(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus) *remove_allowed_grant_type.Service {
	t.Helper()
	finder := domain_services.NewClientFinder(repo)
	return remove_allowed_grant_type.NewService(finder, repo, bus)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
