package _event_handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/_event_handlers"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func aPRForEventHandlers(t *testing.T, userId ids2.UserId) *protected_resources.ProtectedResource {
	t.Helper()
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "A Resource", protected_resources.NewOwnership(ids2.NewUserId(), userId), false)
	require.NoError(t, err)
	pr.PullDomainEvents()
	return pr
}

func TestDisableProtectedResourcesOnUserDeleted(t *testing.T) {
	t.Parallel()

	t.Run("disables every enabled protected resource of the deleted user", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()

		pr := aPRForEventHandlers(t, userId)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.Anything)

		handler := event_handlers.NewDisableProtectedResourcesOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.False(t, pr.IsEnabled())
	})

	t.Run("skips protected resources that are already disabled", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()

		pr := aPRForEventHandlers(t, userId)
		pr.Disable()
		pr.PullDomainEvents()

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)

		handler := event_handlers.NewDisableProtectedResourcesOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository search failure", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		handler := event_handlers.NewDisableProtectedResourcesOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(ids2.NewUserId().String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository save failure", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		userId := ids2.NewUserId()

		pr := aPRForEventHandlers(t, userId)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(fmt.Errorf("db error"))

		handler := event_handlers.NewDisableProtectedResourcesOnUserDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewUserDeleted(userId.String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})
}
