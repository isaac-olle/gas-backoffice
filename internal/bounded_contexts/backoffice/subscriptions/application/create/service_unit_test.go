package create

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/stretchr/testify/require"
)

type createSubscriptionRepository struct {
	saved *subscriptions.Subscription
	err   error
}

func (r *createSubscriptionRepository) Save(_ context.Context, s *subscriptions.Subscription) error {
	r.saved = s
	return r.err
}
func (r *createSubscriptionRepository) SearchById(context.Context, ids.SubscriptionId) (*subscriptions.Subscription, error) {
	return nil, nil
}
func (r *createSubscriptionRepository) List(context.Context, []ids.SubscriptionId, bool) ([]subscriptions.Summary, error) {
	return nil, nil
}

type createSubscriptionBus struct{ events int }

func (b *createSubscriptionBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestCreate(t *testing.T) {
	id, owner := ids.NewSubscriptionId(), ids.NewUserId()
	tests := []struct {
		name, subscriptionName string
		saveErr                error
		wantErr                bool
		wantEvents             int
	}{
		{name: "creates a subscription", subscriptionName: "Team", wantEvents: 1},
		{name: "rejects an empty name", wantErr: true},
		{name: "maps concurrent modification errors", subscriptionName: "Team", saveErr: persistance.ErrConcurrentModification, wantErr: true},
		{name: "wraps repository errors", subscriptionName: "Team", saveErr: errors.New("write failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &createSubscriptionRepository{err: tt.saveErr}, &createSubscriptionBus{}
			err := NewService(repo, bus).Create(context.Background(), NewCommand(id, tt.subscriptionName, owner, false))
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantEvents, bus.events)
			if !tt.wantErr {
				require.NotNil(t, repo.saved)
			}
		})
	}
}
