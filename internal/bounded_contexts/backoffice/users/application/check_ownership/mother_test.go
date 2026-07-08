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
	userId  string
	adminId ids2.AdminId
}

func withUserId(userId string) queryOpt {
	return func(b *queryBuilder) { b.userId = userId }
}

func withAdminId(adminId ids2.AdminId) queryOpt {
	return func(b *queryBuilder) { b.adminId = adminId }
}

// aQuery builds a check_ownership.Query using the Object Mother pattern.
func aQuery(t *testing.T, opts ...queryOpt) *check_ownership.Query {
	t.Helper()
	adminId, err := ids2.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)
	b := &queryBuilder{
		userId:  uuid.New().String(),
		adminId: adminId,
	}
	for _, o := range opts {
		o(b)
	}
	return check_ownership.NewQuery(b.userId, b.adminId)
}

// aUser builds a User aggregate owned by the given admin.
func aUser(t *testing.T, userId ids2.UserId, adminId ids2.AdminId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(adminId))
	require.NoError(t, err)
	user.PullDomainEvents() // clear creation event
	return user
}
