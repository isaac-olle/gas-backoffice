package soft_delete

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

type softDeleteRepository struct {
	user               *users.User
	findErr, deleteErr error
	deleted            ids.UserId
}

func (r *softDeleteRepository) Save(context.Context, *users.User) error { return nil }
func (r *softDeleteRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return r.user, r.findErr
}
func (r *softDeleteRepository) Matching(context.Context, criteria.Criteria) ([]*users.User, error) {
	return nil, nil
}
func (r *softDeleteRepository) SoftDelete(_ context.Context, id ids.UserId) error {
	r.deleted = id
	return r.deleteErr
}
func (r *softDeleteRepository) HardDelete(context.Context, ids.UserId) error { return nil }

type softDeleteBus struct{ events int }

func (b *softDeleteBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestSoftDelete(t *testing.T) {
	id, adminID := ids.NewUserId(), ids.NewUserId()
	author, err := vo.NewUserAuthorFromString(adminID.String(), adminID.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	regular, err := users.NewUser(id, "Ada", "ada@example.com", vo.SubscriptionUserRole, false)
	require.NoError(t, err)
	system, err := users.NewUser(id, "System", "system@example.com", vo.SystemAdministratorRole, true)
	require.NoError(t, err)
	tests := []struct {
		name               string
		user               *users.User
		findErr, deleteErr error
		wantErr            bool
	}{
		{name: "deletes an accessible user", user: regular},
		{name: "rejects a missing user", wantErr: true},
		{name: "wraps lookup errors", findErr: errors.New("lookup failed"), wantErr: true},
		{name: "rejects a system user", user: system, wantErr: true},
		{name: "reports persistence failures", user: regular, deleteErr: errors.New("delete failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &softDeleteRepository{user: tt.user, findErr: tt.findErr, deleteErr: tt.deleteErr}, &softDeleteBus{}
			err := NewService(domain_services.NewUserFinder(repo), repo, bus).SoftDelete(context.Background(), NewCommand(id, author))
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
