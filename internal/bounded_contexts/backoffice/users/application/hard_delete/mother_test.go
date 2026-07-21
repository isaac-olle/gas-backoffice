package hard_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/hard_delete"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	userId string
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	b := &commandBuilder{userId: uuid.New().String()}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.userId)
}

func aUser(t *testing.T) *users.User {
	t.Helper()
	userId := ids2.NewUserId()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(userId))
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}

func aUserWithId(t *testing.T, userId ids2.UserId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(userId))
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}
