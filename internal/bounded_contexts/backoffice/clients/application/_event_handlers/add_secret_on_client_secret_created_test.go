package _event_handlers_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestAddSecretOnClientSecretCreated(t *testing.T) {
	t.Parallel()

	t.Run("adds secret id to client and publishes event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		c, clientId := aClient(t)
		secretId := ids2.NewClientSecretId()
		evt := securityevents.NewClientSecretCreated(secretId.String(), clientId.String())

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
				return len(evts) == 0
			}))

		handler := event_handlers.NewAddSecretOnClientSecretCreatedEventHandler(
			domain_services.NewClientFinder(repo), repo, bus,
		)
		err := handler.Handle(context.Background(), evt)
		require.NoError(t, err)
		assert.Contains(t, saved.ToPrimitives().SecretIds, secretId.String())
	})

	errCases := []struct {
		desc      string
		evt       securityevents.ClientSecretCreated
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id in event",
			evt:      securityevents.NewClientSecretCreated(ids2.NewClientSecretId().String(), "not-a-uuid"),
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns not found when client does not exist",
			evt:  securityevents.NewClientSecretCreated(ids2.NewClientSecretId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "wraps repository save failure",
			evt:  securityevents.NewClientSecretCreated(ids2.NewClientSecretId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewAdminId(), ids2.NewTenantId()), nil, nil)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
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
			handler := event_handlers.NewAddSecretOnClientSecretCreatedEventHandler(
				domain_services.NewClientFinder(repo), repo, bus,
			)
			err := handler.Handle(context.Background(), tc.evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}
