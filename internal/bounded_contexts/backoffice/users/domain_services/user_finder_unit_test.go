package domain_services_test

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

type userFinderRepository struct {
	user *users.User
	err  error
}

func (r *userFinderRepository) Save(context.Context, *users.User) error { return nil }
func (r *userFinderRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return r.user, r.err
}
func (r *userFinderRepository) Matching(context.Context, criteria.Criteria) ([]*users.User, error) {
	return nil, nil
}
func (r *userFinderRepository) SoftDelete(context.Context, ids.UserId) error { return nil }
func (r *userFinderRepository) HardDelete(context.Context, ids.UserId) error { return nil }

func TestUserFinderFind(t *testing.T) {
	id := ids.NewUserId()
	user, err := users.NewUser(id, "Ada Lovelace", "ada@example.com", vo.SystemAdministratorRole, true)
	require.NoError(t, err)
	tests := []struct {
		name    string
		user    *users.User
		repoErr error
		wantErr bool
	}{
		{name: "returns the user", user: user},
		{name: "returns not found for nil result", wantErr: true},
		{name: "wraps repository failures", repoErr: errors.New("database unavailable"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain_services.NewUserFinder(&userFinderRepository{user: tt.user, err: tt.repoErr}).Find(context.Background(), id)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				require.Same(t, user, got)
			}
		})
	}
}
