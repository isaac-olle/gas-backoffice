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

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, prId ids2.ProtectedResourceId) *scopes.Scope
		evtPrId string
		wantErr bool
		check   func(t *testing.T, scope *scopes.Scope)
	}{
		{
			name: "disables every enabled scope of the deleted protected resource",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, prId ids2.ProtectedResourceId) *scopes.Scope {
				scope := aScopeForEventHandlers(t, prId, ids2.NewSubscriptionId())
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{scope}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(nil)
				bus.EXPECT().Publish(mock.Anything, mock.Anything)
				return scope
			},
			check: func(t *testing.T, scope *scopes.Scope) {
				assert.False(t, scope.IsEnabled())
			},
		},
		{
			name: "skips scopes that are already disabled",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, prId ids2.ProtectedResourceId) *scopes.Scope {
				scope := aScopeForEventHandlers(t, prId, ids2.NewSubscriptionId())
				scope.Disable()
				scope.PullDomainEvents()
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{scope}, nil)
				return scope
			},
		},
		{
			name:    "rejects invalid protected resource id in event",
			evtPrId: "not-a-uuid",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, prId ids2.ProtectedResourceId) *scopes.Scope {
				return nil
			},
			wantErr: true,
		},
		{
			name: "wraps repository search failure",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, prId ids2.ProtectedResourceId) *scopes.Scope {
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
				return nil
			},
			wantErr: true,
		},
		{
			name: "wraps repository save failure",
			setup: func(t *testing.T, repo *scopemocks.MockRepository, bus *eventmocks.MockBus, prId ids2.ProtectedResourceId) *scopes.Scope {
				scope := aScopeForEventHandlers(t, prId, ids2.NewSubscriptionId())
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*scopes.Scope{scope}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).Return(fmt.Errorf("db error"))
				return scope
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			prId := ids2.NewProtectedResourceId()

			scope := tt.setup(t, repo, bus, prId)

			evtPrId := tt.evtPrId
			if evtPrId == "" {
				evtPrId = prId.String()
			}
			handler := event_handlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(repo, bus)
			evt := backofficeevents.NewProtectedResourceDeleted(evtPrId, ids2.NewUserId().String())
			err := handler.Handle(context.Background(), evt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, scope)
			}
		})
	}
}
