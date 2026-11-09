package hard_delete

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type hardDeleteRepository struct {
	user           *users.User
	err, deleteErr error
	deleted        ids.UserId
}

func (r *hardDeleteRepository) Save(context.Context, *users.User) error { return nil }
func (r *hardDeleteRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return r.user, r.err
}
func (r *hardDeleteRepository) Matching(context.Context, criteria.Criteria) ([]*users.User, error) {
	return nil, nil
}
func (r *hardDeleteRepository) SoftDelete(context.Context, ids.UserId) error { return nil }
func (r *hardDeleteRepository) HardDelete(_ context.Context, id ids.UserId) error {
	r.deleted = id
	return r.deleteErr
}

type hardDeleteBus struct{ events int }

func (b *hardDeleteBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestHardDelete(t *testing.T) {
	id, adminID := ids.NewUserId(), ids.NewUserId()
	admin, err := vo.NewUserAuthorFromString(adminID.String(), adminID.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	regular, err := users.NewUser(id, "Ada", "ada@example.com", vo.SubscriptionUserRole, false)
	require.NoError(t, err)
	system, err := users.NewUser(id, "System", "system@example.com", vo.SystemAdministratorRole, true)
	require.NoError(t, err)
	tests := []struct {
		name, rawID        string
		user               *users.User
		findErr, deleteErr error
		wantErr            bool
	}{
		{name: "deletes an accessible user", rawID: id.String(), user: regular},
		{name: "rejects an invalid id", rawID: "bad-id", wantErr: true},
		{name: "rejects a missing user", rawID: id.String(), wantErr: true},
		{name: "rejects deleting a system user", rawID: id.String(), user: system, wantErr: true},
		{name: "reports persistence failures", rawID: id.String(), user: regular, deleteErr: errors.New("delete failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &hardDeleteRepository{user: tt.user, err: tt.findErr, deleteErr: tt.deleteErr}, &hardDeleteBus{}
			err := NewService(domain_services.NewUserFinder(repo), repo, bus).HardDelete(context.Background(), NewCommand(tt.rawID, admin))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				require.Zero(t, bus.events)
				return
			}
			require.Equal(t, id, repo.deleted)
			require.Equal(t, 2, bus.events)
		})
	}
}
