package change_name

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
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/stretchr/testify/require"
)

type changeNameRepository struct {
	user             *users.User
	findErr, saveErr error
	saves            int
}

func (r *changeNameRepository) Save(context.Context, *users.User) error { r.saves++; return r.saveErr }
func (r *changeNameRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return r.user, r.findErr
}
func (r *changeNameRepository) Matching(context.Context, criteria.Criteria) ([]*users.User, error) {
	return nil, nil
}
func (r *changeNameRepository) SoftDelete(context.Context, ids.UserId) error { return nil }
func (r *changeNameRepository) HardDelete(context.Context, ids.UserId) error { return nil }

type changeNameBus struct{ events int }

func (b *changeNameBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestChangeName(t *testing.T) {
	id, adminID, otherID := ids.NewUserId(), ids.NewUserId(), ids.NewUserId()
	admin, err := vo.NewUserAuthorFromString(adminID.String(), adminID.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	other, err := vo.NewUserAuthorFromString(otherID.String(), otherID.String(), vo.SubscriptionUserRole)
	require.NoError(t, err)
	regular, err := users.NewUser(id, "Ada", "ada@example.com", vo.SubscriptionUserRole, false)
	require.NoError(t, err)
	system, err := users.NewUser(id, "System", "system@example.com", vo.SystemAdministratorRole, true)
	require.NoError(t, err)
	tests := []struct {
		name, newName    string
		user             *users.User
		author           *vo.UserAuthor
		findErr, saveErr error
		wantErr          bool
	}{
		{name: "changes the name", newName: "Ada Lovelace", user: regular, author: admin},
		{name: "rejects an invalid name", user: regular, author: admin, wantErr: true},
		{name: "wraps finder errors", newName: "Ada", findErr: errors.New("lookup failed"), author: admin, wantErr: true},
		{name: "rejects callers without access", newName: "Ada", user: regular, author: other, wantErr: true},
		{name: "rejects changes to system users", newName: "System 2", user: system, author: admin, wantErr: true},
		{name: "wraps ordinary persistence failures", newName: "Ada 2", user: regular, author: admin, saveErr: errors.New("write failed"), wantErr: true},
		{name: "maps concurrent modifications to conflict", newName: "Ada 2", user: regular, author: admin, saveErr: persistance.ErrConcurrentModification, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &changeNameRepository{user: tt.user, findErr: tt.findErr, saveErr: tt.saveErr}, &changeNameBus{}
			err := NewService(domain_services.NewUserFinder(repo), repo, bus).ChangeName(context.Background(), NewCommand(id, tt.newName, tt.author))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				require.Zero(t, bus.events)
				return
			}
			require.Equal(t, 1, repo.saves)
			require.Equal(t, 2, bus.events)
		})
	}
}
