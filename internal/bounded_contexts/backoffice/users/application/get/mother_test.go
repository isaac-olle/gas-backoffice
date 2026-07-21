package get_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/stretchr/testify/require"
)

func aUserWithId(t *testing.T, userId ids2.UserId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(userId))
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}
