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

func aClientForEventHandlers(t *testing.T, userId ids2.UserId) *clients.Client {
	t.Helper()
	client, err := clients.NewClient(ids2.NewClientId(), "A Client", clients.NewOwnership(userId, userId), nil, nil, false, false)
	require.NoError(t, err)
	client.PullDomainEvents()
	return client
}

func TestDisableClientsOnUserDeleted(t *testing.T) {
	t.Parallel()

	t.Run("disables every enabled client of the deleted user", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()

		client := aClientForEventHandlers(t, userId)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{client}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.Anything)

		handler := event_handlers.NewDisableClientsOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.False(t, client.IsEnabled())
	})

	t.Run("skips clients that are already disabled", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()

		client := aClientForEventHandlers(t, userId)
		client.Disable()
		client.PullDomainEvents()

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{client}, nil)

		handler := event_handlers.NewDisableClientsOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository search failure", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		handler := event_handlers.NewDisableClientsOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(ids2.NewUserId().String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository save failure", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()

		client := aClientForEventHandlers(t, userId)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*clients.Client{client}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))

		handler := event_handlers.NewDisableClientsOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(userId.String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})
}
