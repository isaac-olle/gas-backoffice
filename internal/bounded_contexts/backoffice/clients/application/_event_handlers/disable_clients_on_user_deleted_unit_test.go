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
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func aClientForEventHandlers(t *testing.T, subscriptionId ids2.SubscriptionId) *clients.Client {
	t.Helper()
	client, err := clients.NewClient(ids2.NewClientId(), "A Client", clients.NewOwnership(subscriptionId), nil, nil, false, false)
	require.NoError(t, err)
	client.PullDomainEvents()
	return client
}

func TestDisableClientsOnUserDeleted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus, userId ids2.UserId) *clients.Client
		wantErr bool
		check   func(t *testing.T, client *clients.Client)
	}{
		{
			name: "disables every enabled client of the deleted user",
			setup: func(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus, userId ids2.UserId) *clients.Client {
				client := aClientForEventHandlers(t, ids2.NewSubscriptionId())
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{client}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(nil)
				bus.EXPECT().Publish(mock.Anything, mock.Anything)
				return client
			},
			check: func(t *testing.T, client *clients.Client) {
				assert.False(t, client.IsEnabled())
			},
		},
		{
			name: "skips clients that are already disabled",
			setup: func(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus, userId ids2.UserId) *clients.Client {
				client := aClientForEventHandlers(t, ids2.NewSubscriptionId())
				client.Disable()
				client.PullDomainEvents()
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{client}, nil)
				return client
			},
		},
		{
			name: "wraps repository search failure",
			setup: func(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus, userId ids2.UserId) *clients.Client {
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
				return nil
			},
			wantErr: true,
		},
		{
			name: "wraps repository save failure",
			setup: func(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus, userId ids2.UserId) *clients.Client {
				client := aClientForEventHandlers(t, ids2.NewSubscriptionId())
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{client}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
				return client
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := clientmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			userId := ids2.NewUserId()

			client := tt.setup(t, repo, bus, userId)

			handler := event_handlers.NewDisableClientsOnUserDeletedEventHandler(repo, bus)
			evt := backofficeevents.NewUserDeleted(userId.String())
			err := handler.Handle(context.Background(), evt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, client)
			}
		})
	}
}
