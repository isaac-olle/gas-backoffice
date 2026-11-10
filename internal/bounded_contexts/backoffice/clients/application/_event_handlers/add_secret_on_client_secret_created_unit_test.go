package _event_handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

func TestAddSecretOnClientSecretCreated(t *testing.T) {
	t.Parallel()

	newHandler := func(repo *clientmocks.MockRepository) *event_handlers.AddSecretOnClientSecretCreatedEventHandler {
		return event_handlers.NewAddSecretOnClientSecretCreatedEventHandler(domain_services.NewClientFinder(repo), repo)
	}

	t.Run("adds secret id to client and saves it", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		client := aClientForEventHandlers(t, ids2.NewSubscriptionId())
		clientId, err := ids2.NewClientIdFromString(client.ToPrimitives().Id)
		require.NoError(t, err)
		childId := ids2.NewClientSecretId()

		evt := securityevents.NewClientSecretCreated(childId.String(), clientId.String())

		repo.EXPECT().SearchById(mock.Anything, clientId).Return(client, nil)
		repo.EXPECT().Save(mock.Anything, client).Return(nil)

		require.NoError(t, newHandler(repo).Handle(context.Background(), evt))
		assert.Contains(t, client.ToPrimitives().SecretIds, childId.String())
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
			desc:     "rejects invalid secret id in event",
			evt:      securityevents.NewClientSecretCreated("not-a-uuid", ids2.NewClientId().String()),
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "fails when client does not exist",
			evt:  securityevents.NewClientSecretCreated(ids2.NewClientSecretId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "maps concurrent modification to conflict",
			evt:  securityevents.NewClientSecretCreated(ids2.NewClientSecretId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(aClientForEventHandlers(t, ids2.NewSubscriptionId()), nil)
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(persistance.ErrConcurrentModification)
			},
			wantCode: clients.ErrClientConcurrentModification,
		},
		{
			desc: "wraps repository save failure",
			evt:  securityevents.NewClientSecretCreated(ids2.NewClientSecretId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(aClientForEventHandlers(t, ids2.NewSubscriptionId()), nil)
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			err := newHandler(repo).Handle(context.Background(), tc.evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}
