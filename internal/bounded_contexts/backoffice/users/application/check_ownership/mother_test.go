package check_ownership_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_ownership"
	"github.com/stretchr/testify/require"
)

// queryOpt is a functional option for queryBuilder.
type queryOpt func(*queryBuilder)

type queryBuilder struct {
	userId string
}

func withUserId(userId string) queryOpt {
	return func(b *queryBuilder) { b.userId = userId }
}

// aQuery builds a check_ownership.Query using the Object Mother pattern.
func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	b := &queryBuilder{
		userId: uuid.New().String(),
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.userId)
}

// aUser builds a User aggregate with the given userId as ownership.
func aUser(t *testing.T, userId ids2.UserId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(userId))
	require.NoError(t, err)
	user.PullDomainEvents() // clear creation event
	return user
}
