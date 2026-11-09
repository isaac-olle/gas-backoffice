package _event_handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/_event_handlers"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func TestDisableScopesOnProtectedResourceDeleted(t *testing.T) {
	t.Parallel()

	t.Run("disables every enabled scope of the deleted protected resource", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		scope := aScopeForEventHandlers(t, prId, userId)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{scope}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.Anything)

		handler := event_handlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewProtectedResourceDeleted(prId.String(), userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.False(t, scope.IsEnabled())
	})

	t.Run("skips scopes that are already disabled", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		scope := aScopeForEventHandlers(t, prId, userId)
		scope.Disable()
		scope.PullDomainEvents()

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{scope}, nil)

		handler := event_handlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewProtectedResourceDeleted(prId.String(), userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))
	})

	t.Run("rejects invalid protected resource id in event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		handler := event_handlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewProtectedResourceDeleted("not-a-uuid", ids2.NewUserId().String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository search failure", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		handler := event_handlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewProtectedResourceDeleted(ids2.NewProtectedResourceId().String(), ids2.NewUserId().String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository save failure", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		scope := aScopeForEventHandlers(t, prId, userId)

		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{scope}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(fmt.Errorf("db error"))

		handler := event_handlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewProtectedResourceDeleted(prId.String(), userId.String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})
}
