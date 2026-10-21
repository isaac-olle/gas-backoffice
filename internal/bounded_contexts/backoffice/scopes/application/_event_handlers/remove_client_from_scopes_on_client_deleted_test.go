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

func aScopeForEventHandlers(t *testing.T, prId ids2.ProtectedResourceId, userId ids2.UserId) *scopes.Scope {
	t.Helper()
	scope, err := scopes.NewScope(ids2.NewScopeId(), "read:things", scopes.NewOwnershipFromVos(prId, ids2.NewUserId(), userId), false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}

func TestRemoveClientFromScopesOnClientDeleted(t *testing.T) {
	t.Parallel()

	t.Run("removes the deleted client from every affected scope", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		clientId := ids2.NewClientId()
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		scope := aScopeForEventHandlers(t, prId, userId)
		scope.AddAllowedClient(clientId)
		scope.PullDomainEvents()

		repo.EXPECT().FindByAllowedClientId(mock.Anything, clientId).Return([]*scopes.Scope{scope}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.Anything)

		handler := event_handlers.NewRemoveClientFromScopesOnClientDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewClientDeleted(clientId.String(), userId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.NotContains(t, scope.ToPrimitives().AllowedClientIds, clientId.String())
	})

	t.Run("rejects invalid client id in event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		handler := event_handlers.NewRemoveClientFromScopesOnClientDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewClientDeleted("not-a-uuid", ids2.NewUserId().String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository search failure", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().FindByAllowedClientId(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

		handler := event_handlers.NewRemoveClientFromScopesOnClientDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewClientDeleted(ids2.NewClientId().String(), ids2.NewUserId().String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})

	t.Run("wraps repository save failure", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		clientId := ids2.NewClientId()
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		scope := aScopeForEventHandlers(t, prId, userId)
		scope.AddAllowedClient(clientId)
		scope.PullDomainEvents()

		repo.EXPECT().FindByAllowedClientId(mock.Anything, clientId).Return([]*scopes.Scope{scope}, nil)
		repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(fmt.Errorf("db error"))

		handler := event_handlers.NewRemoveClientFromScopesOnClientDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewClientDeleted(clientId.String(), userId.String())
		assert.Error(t, handler.Handle(context.Background(), evt))
	})
}
