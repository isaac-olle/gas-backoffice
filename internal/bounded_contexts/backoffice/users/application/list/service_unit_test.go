package list

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type listUsersRepository struct {
	result []*users.User
	err    error
	query  criteria.Criteria
}

func (r *listUsersRepository) Save(context.Context, *users.User) error { return nil }
func (r *listUsersRepository) SearchById(context.Context, ids.UserId) (*users.User, error) {
	return nil, nil
}
func (r *listUsersRepository) Matching(_ context.Context, query criteria.Criteria) ([]*users.User, error) {
	r.query = query
	return r.result, r.err
}
func (r *listUsersRepository) SoftDelete(context.Context, ids.UserId) error { return nil }
func (r *listUsersRepository) HardDelete(context.Context, ids.UserId) error { return nil }

func TestList(t *testing.T) {
	author := ids.NewUserId()
	rows := []*users.User{}
	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "returns matching users"},
		{name: "wraps repository errors", repoErr: errors.New("database unavailable"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &listUsersRepository{result: rows, err: tt.repoErr}
			got, err := NewService(repo).List(context.Background(), NewQuery(author))
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			require.Equal(t, rows, got)
			require.Equal(t, "author_id", repo.query.Filters[0].Field)
			require.Equal(t, author.String(), repo.query.Filters[0].Value)
		})
	}
}
