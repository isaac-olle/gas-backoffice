package list

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type listRepositoryStub struct {
	result []subscriptions.Summary
	err    error
	gotIDs []ids.SubscriptionId
	gotAll bool
}

func (r *listRepositoryStub) Save(context.Context, *subscriptions.Subscription) error { return nil }
func (r *listRepositoryStub) SearchById(context.Context, ids.SubscriptionId) (*subscriptions.Subscription, error) {
	return nil, nil
}
func (r *listRepositoryStub) List(_ context.Context, ids []ids.SubscriptionId, all bool) ([]subscriptions.Summary, error) {
	r.gotIDs, r.gotAll = ids, all
	return r.result, r.err
}

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		all     bool
		ids     []ids.SubscriptionId
		result  []subscriptions.Summary
		repoErr error
		wantErr bool
	}{
		{name: "lists all subscriptions", all: true, result: []subscriptions.Summary{{Id: "one"}}},
		{name: "lists selected subscriptions", ids: []ids.SubscriptionId{ids.NewSubscriptionId()}, result: []subscriptions.Summary{{Id: "selected"}}},
		{name: "wraps repository errors", repoErr: errors.New("query failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &listRepositoryStub{result: tt.result, err: tt.repoErr}
			got, err := NewService(repo).List(context.Background(), NewQuery(tt.all, tt.ids))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.result, got)
			require.Equal(t, tt.ids, repo.gotIDs)
			require.Equal(t, tt.all, repo.gotAll)
		})
	}
}
