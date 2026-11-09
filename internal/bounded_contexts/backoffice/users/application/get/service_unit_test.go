package get

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type getUserRepository struct {
	user *users.User
	err  error
}

func (r *getUserRepository) Save(context.Context, *users.User) error { return nil }
func (r *getUserRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return r.user, r.err
}
func (r *getUserRepository) Matching(context.Context, criteria.Criteria) ([]*users.User, error) {
	return nil, nil
}
func (r *getUserRepository) SoftDelete(context.Context, ids.UserId) error { return nil }
func (r *getUserRepository) HardDelete(context.Context, ids.UserId) error { return nil }

func TestGet(t *testing.T) {
	userID, adminID := ids.NewUserId(), ids.NewUserId()
	user, err := users.NewUser(userID, "Ada Lovelace", "ada@example.com", vo.SubscriptionUserRole, false)
	require.NoError(t, err)
	author, err := vo.NewUserAuthorFromString(adminID.String(), adminID.String(), vo.SystemAdministratorRole)
	require.NoError(t, err)
	otherId := ids.NewUserId()
	other, err := vo.NewUserAuthorFromString(otherId.String(), otherId.String(), vo.SubscriptionUserRole)
	require.NoError(t, err)
	tests := []struct {
		name    string
		found   *users.User
		repoErr error
		caller  *vo.UserAuthor
		wantErr bool
	}{
		{name: "returns an accessible user", found: user, caller: author},
		{name: "returns not found when the finder has no user", caller: author, wantErr: true},
		{name: "wraps finder errors", repoErr: errors.New("lookup failed"), caller: author, wantErr: true},
		{name: "rejects callers without access", found: user, caller: other, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(domain_services.NewUserFinder(&getUserRepository{user: tt.found, err: tt.repoErr})).Get(context.Background(), NewQuery(userID, tt.caller))
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				require.Same(t, user, got)
			}
		})
	}
}
