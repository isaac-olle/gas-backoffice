package create_user

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type createUserRepository struct {
	saved *users.User
	err   error
}

func (r *createUserRepository) Save(_ context.Context, u *users.User) error {
	r.saved = u
	return r.err
}
func (r *createUserRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return nil, nil
}
func (r *createUserRepository) Matching(context.Context, criteria.Criteria) ([]*users.User, error) {
	return nil, nil
}
func (r *createUserRepository) SoftDelete(context.Context, ids.UserId) error { return nil }
func (r *createUserRepository) HardDelete(context.Context, ids.UserId) error { return nil }

type createUserBus struct{ published int }

func (b *createUserBus) Publish(context.Context, []event.Event) { b.published++ }

func TestCreate(t *testing.T) {
	tests := []struct {
		name, userName, email string
		repoErr               error
		wantErr               bool
	}{
		{name: "creates and publishes a valid user", userName: "Ada Lovelace", email: "ada@example.com"},
		{name: "rejects an invalid name", userName: "", email: "ada@example.com", wantErr: true},
		{name: "rejects an invalid email", userName: "Ada Lovelace", email: "invalid", wantErr: true},
		{name: "reports save errors", userName: "Ada Lovelace", email: "ada@example.com", repoErr: errors.New("write failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, bus := &createUserRepository{err: tt.repoErr}, &createUserBus{}
			cmd := NewCommand(ids.NewUserId(), ids.NewUserId(), tt.userName, vo.SystemAdministratorRole, false, tt.email)
			err := NewService(repo, bus).Create(context.Background(), cmd)
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				require.Zero(t, bus.published)
				return
			}
			require.NotNil(t, repo.saved)
			require.Equal(t, 1, bus.published)
		})
	}
}
