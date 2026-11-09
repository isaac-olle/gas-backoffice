package add_member

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type addMemberRepository struct {
	aggregate        *subscriptions.Subscription
	findErr, saveErr error
	saves            int
}

func (r *addMemberRepository) Save(_ context.Context, s *subscriptions.Subscription) error {
	r.saves++
	r.aggregate = s
	return r.saveErr
}
func (r *addMemberRepository) SearchById(context.Context, ids.SubscriptionId) (*subscriptions.Subscription, error) {
	return r.aggregate, r.findErr
}
func (r *addMemberRepository) List(context.Context, []ids.SubscriptionId, bool) ([]subscriptions.Summary, error) {
	return nil, nil
}

type addMemberBus struct{ calls, events int }

func (b *addMemberBus) Publish(_ context.Context, events []event.Event) {
	b.calls++
	b.events += len(events)
}

func TestAdd(t *testing.T) {
	subID, ownerID, memberID := ids.NewSubscriptionId(), ids.NewUserId(), ids.NewUserId()
	base, err := subscriptions.NewSubscription(subID, "Team", ownerID, false)
	require.NoError(t, err)
	duplicate, err := subscriptions.FromPrimitives(subscriptions.SubscriptionPrimitives{Id: subID.String(), Name: "Team", OwnerId: ownerID.String(), Members: []subscriptions.MembershipPrimitives{{UserId: memberID.String(), SubscriptionRole: string(vo.UserSubscriptionRole)}}})
	require.NoError(t, err)
	tests := []struct {
		name             string
		aggregate        *subscriptions.Subscription
		findErr, saveErr error
		wantErr          bool
		wantEvents       int
	}{
		{name: "adds a member", aggregate: base, wantEvents: 2},
		{name: "rejects a missing subscription", wantErr: true},
		{name: "wraps lookup errors", findErr: errors.New("lookup failed"), wantErr: true},
		{name: "wraps save errors", aggregate: base, saveErr: errors.New("save failed"), wantErr: true},
		{name: "does not record a duplicate member addition", aggregate: duplicate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &addMemberRepository{aggregate: tt.aggregate, findErr: tt.findErr, saveErr: tt.saveErr}, &addMemberBus{}
			err := NewService(repo, bus).Add(context.Background(), NewCommand(subID, memberID, vo.UserSubscriptionRole))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				require.Zero(t, bus.calls)
				return
			}
			require.Equal(t, 1, repo.saves)
			require.Equal(t, 1, bus.calls)
			require.Equal(t, tt.wantEvents, bus.events)
			require.True(t, repo.aggregate.HasMember(memberID))
		})
	}
}
