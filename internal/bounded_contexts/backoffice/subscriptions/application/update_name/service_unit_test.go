package update_name

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

type updateSubscriptionRepository struct {
	aggregate        *subscriptions.Subscription
	findErr, saveErr error
	saved            *subscriptions.Subscription
}

func (r *updateSubscriptionRepository) Save(_ context.Context, s *subscriptions.Subscription) error {
	r.saved = s
	return r.saveErr
}
func (r *updateSubscriptionRepository) SearchById(context.Context, ids.SubscriptionId) (*subscriptions.Subscription, error) {
	return r.aggregate, r.findErr
}
func (r *updateSubscriptionRepository) List(context.Context, []ids.SubscriptionId, bool) ([]subscriptions.Summary, error) {
	return nil, nil
}

type updateSubscriptionBus struct{ events int }

func (b *updateSubscriptionBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestUpdateName(t *testing.T) {
	subID, owner := ids.NewSubscriptionId(), ids.NewUserId()
	regular, err := subscriptions.NewSubscription(subID, "Old", owner, false)
	require.NoError(t, err)
	system, err := subscriptions.NewSubscription(subID, "System", owner, true)
	require.NoError(t, err)
	tests := []struct {
		name, newName    string
		aggregate        *subscriptions.Subscription
		findErr, saveErr error
		wantErr          bool
		wantEvents       int
	}{
		{name: "updates an existing subscription", newName: "New", aggregate: regular, wantEvents: 1},
		{name: "creates a subscription when the id is new", newName: "New", wantEvents: 1},
		{name: "rejects an invalid name during upsert", aggregate: nil, wantErr: true},
		{name: "rejects an invalid name on an existing subscription", aggregate: regular, wantErr: true},
		{name: "rejects changes to a system subscription", newName: "New", aggregate: system, wantErr: true},
		{name: "wraps lookup failures", newName: "New", findErr: errors.New("lookup failed"), wantErr: true},
		{name: "maps concurrent modification errors", newName: "New", aggregate: regular, saveErr: persistance.ErrConcurrentModification, wantErr: true},
		{name: "wraps other save failures", newName: "New", aggregate: regular, saveErr: errors.New("write failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &updateSubscriptionRepository{aggregate: tt.aggregate, findErr: tt.findErr, saveErr: tt.saveErr}, &updateSubscriptionBus{}
			err := NewService(repo, bus).UpdateName(context.Background(), NewCommand(subID, tt.newName, owner))
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantEvents, bus.events)
			if !tt.wantErr {
				require.NotNil(t, repo.saved)
				require.Equal(t, "New", repo.saved.ToPrimitives().Name)
			}
		})
	}
}
