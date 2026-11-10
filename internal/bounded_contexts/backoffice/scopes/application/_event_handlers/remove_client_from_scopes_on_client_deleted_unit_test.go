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

func aScopeForEventHandlers(t *testing.T, prId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) *scopes.Scope {
	t.Helper()
	scope, err := scopes.NewScope(ids2.NewScopeId(), "read:things", scopes.NewOwnershipFromVos(prId, subscriptionId), false)
	require.NoError(t, err)
	scope.PullDomainEvents()
	return scope
}

func TestRemoveClientFromScopesOnClientDeleted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, clientId ids2.ClientId) (backofficeevents.ClientDeleted, *scopes.Scope)
		wantErr bool
	}{
		{
			name: "removes the deleted client from every affected scope",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, clientId ids2.ClientId) (backofficeevents.ClientDeleted, *scopes.Scope) {
				prId := ids2.NewProtectedResourceId()
				subscriptionId := ids2.NewSubscriptionId()
				scope := aScopeForEventHandlers(t, prId, subscriptionId)
				scope.AddAllowedClient(clientId)
				scope.PullDomainEvents()

				repo.EXPECT().FindByAllowedClientId(mock.Anything, clientId).Return([]*scopes.Scope{scope}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(nil)
				bus.EXPECT().Publish(mock.Anything, mock.Anything)

				return backofficeevents.NewClientDeleted(clientId.String(), ids2.NewUserId().String()), scope
			},
		},
		{
			name: "rejects invalid client id in event",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, clientId ids2.ClientId) (backofficeevents.ClientDeleted, *scopes.Scope) {
				return backofficeevents.NewClientDeleted("not-a-uuid", ids2.NewUserId().String()), nil
			},
			wantErr: true,
		},
		{
			name: "wraps repository search failure",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, clientId ids2.ClientId) (backofficeevents.ClientDeleted, *scopes.Scope) {
				repo.EXPECT().FindByAllowedClientId(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
				return backofficeevents.NewClientDeleted(ids2.NewClientId().String(), ids2.NewUserId().String()), nil
			},
			wantErr: true,
		},
		{
			name: "wraps repository save failure",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, clientId ids2.ClientId) (backofficeevents.ClientDeleted, *scopes.Scope) {
				prId := ids2.NewProtectedResourceId()
				subscriptionId := ids2.NewSubscriptionId()
				scope := aScopeForEventHandlers(t, prId, subscriptionId)
				scope.AddAllowedClient(clientId)
				scope.PullDomainEvents()

				repo.EXPECT().FindByAllowedClientId(mock.Anything, clientId).Return([]*scopes.Scope{scope}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(fmt.Errorf("db error"))

				return backofficeevents.NewClientDeleted(clientId.String(), ids2.NewUserId().String()), scope
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			clientId := ids2.NewClientId()

			evt, scope := tt.setup(t, repo, bus, clientId)
			err := event_handlers.NewRemoveClientFromScopesOnClientDeletedEventHandler(repo, bus).Handle(context.Background(), evt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotContains(t, scope.ToPrimitives().AllowedClientIds, clientId.String())
		})
	}
}
