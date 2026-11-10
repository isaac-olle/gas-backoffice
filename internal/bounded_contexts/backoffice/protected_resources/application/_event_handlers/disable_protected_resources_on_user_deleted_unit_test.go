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

func aPRForEventHandlers(t *testing.T, subscriptionId ids2.SubscriptionId) *protected_resources.ProtectedResource {
	t.Helper()
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "A Resource", protected_resources.NewOwnership(subscriptionId), false)
	require.NoError(t, err)
	pr.PullDomainEvents()
	return pr
}

func TestDisableProtectedResourcesOnUserDeleted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *prmocks.MockRepository, bus *eventmocks.MockBus) *protected_resources.ProtectedResource
		wantErr bool
		check   func(t *testing.T, pr *protected_resources.ProtectedResource)
	}{
		{
			name: "disables every enabled protected resource of the deleted user",
			setup: func(t *testing.T, repo *prmocks.MockRepository, bus *eventmocks.MockBus) *protected_resources.ProtectedResource {
				pr := aPRForEventHandlers(t, ids2.NewSubscriptionId())
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(nil)
				bus.EXPECT().Publish(mock.Anything, mock.Anything)
				return pr
			},
			check: func(t *testing.T, pr *protected_resources.ProtectedResource) {
				assert.False(t, pr.IsEnabled())
			},
		},
		{
			name: "skips protected resources that are already disabled",
			setup: func(t *testing.T, repo *prmocks.MockRepository, bus *eventmocks.MockBus) *protected_resources.ProtectedResource {
				pr := aPRForEventHandlers(t, ids2.NewSubscriptionId())
				pr.Disable()
				pr.PullDomainEvents()
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)
				return pr
			},
		},
		{
			name: "wraps repository search failure",
			setup: func(t *testing.T, repo *prmocks.MockRepository, bus *eventmocks.MockBus) *protected_resources.ProtectedResource {
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
				return nil
			},
			wantErr: true,
		},
		{
			name: "wraps repository save failure",
			setup: func(t *testing.T, repo *prmocks.MockRepository, bus *eventmocks.MockBus) *protected_resources.ProtectedResource {
				pr := aPRForEventHandlers(t, ids2.NewSubscriptionId())
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(fmt.Errorf("db error"))
				return pr
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := prmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			userId := ids2.NewUserId()

			pr := tt.setup(t, repo, bus)

			handler := event_handlers.NewDisableProtectedResourcesOnUserDeletedEventHandler(repo, bus)
			evt := backofficeevents.NewUserDeleted(userId.String())
			err := handler.Handle(context.Background(), evt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, pr)
			}
		})
	}
}
