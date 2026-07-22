package soft_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/soft_delete"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	userId string
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{userId: uuid.New().String()}
	for _, o := range opts {
		o(b)
	}
	id, err := ids2.NewUserIdFromString(b.userId)
	require.NoError(t, err)
	return soft_delete.NewCommand(id)
}

func aUser(t *testing.T) *users.User {
	t.Helper()
	userId := ids2.NewUserId()
	user, err := users.NewUser(userId, "Test User", "tenant", ids2.NewUserId())
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}

func aUserWithId(t *testing.T, userId ids2.UserId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", "tenant", ids2.NewUserId())
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}
