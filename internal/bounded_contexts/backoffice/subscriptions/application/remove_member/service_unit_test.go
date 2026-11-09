package remove_member

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/stretchr/testify/require"
)

type removeMemberRepository struct {
	aggregate        *subscriptions.Subscription
	findErr, saveErr error
	saves            int
}

func (r *removeMemberRepository) Save(_ context.Context, s *subscriptions.Subscription) error {
	r.saves++
	r.aggregate = s
	return r.saveErr
}
func (r *removeMemberRepository) SearchById(context.Context, ids.SubscriptionId) (*subscriptions.Subscription, error) {
	return r.aggregate, r.findErr
}
func (r *removeMemberRepository) List(context.Context, []ids.SubscriptionId, bool) ([]subscriptions.Summary, error) {
	return nil, nil
}

type removeMemberBus struct{ events int }

func (b *removeMemberBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestRemove(t *testing.T) {
	subID, ownerID, memberID := ids.NewSubscriptionId(), ids.NewUserId(), ids.NewUserId()
	withMember, err := subscriptions.FromPrimitives(subscriptions.SubscriptionPrimitives{Id: subID.String(), Name: "Team", OwnerId: ownerID.String(), Members: []subscriptions.MembershipPrimitives{{UserId: memberID.String(), SubscriptionRole: string(vo.UserSubscriptionRole)}}})
	require.NoError(t, err)
	empty, err := subscriptions.FromPrimitives(subscriptions.SubscriptionPrimitives{Id: subID.String(), Name: "Team", OwnerId: ownerID.String()})
	require.NoError(t, err)
	tests := []struct {
		name             string
		aggregate        *subscriptions.Subscription
		findErr, saveErr error
		wantErr          bool
		wantEvents       int
	}{
		{name: "removes an existing member", aggregate: withMember, wantEvents: 1},
		{name: "is a no-op for a missing member", aggregate: empty},
		{name: "rejects a missing subscription", wantErr: true},
		{name: "wraps lookup errors", findErr: errors.New("lookup failed"), wantErr: true},
		{name: "maps concurrent modification errors", aggregate: withMember, saveErr: persistance.ErrConcurrentModification, wantErr: true},
		{name: "wraps other save errors", aggregate: withMember, saveErr: errors.New("write failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &removeMemberRepository{aggregate: tt.aggregate, findErr: tt.findErr, saveErr: tt.saveErr}, &removeMemberBus{}
			err := NewService(repo, bus).Remove(context.Background(), NewCommand(subID, memberID))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				require.Zero(t, bus.events)
				return
			}
			require.Equal(t, 1, repo.saves)
			require.Equal(t, tt.wantEvents, bus.events)
			require.False(t, repo.aggregate.HasMember(memberID))
		})
	}
}
