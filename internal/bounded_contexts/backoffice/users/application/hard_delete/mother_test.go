package hard_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/hard_delete"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	userId string
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

// aCommand builds a hard_delete.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	b := &commandBuilder{userId: uuid.New().String()}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.userId)
}

// aUser builds a User aggregate for use in tests.
func aUser(t *testing.T) *users.User {
	t.Helper()
	user, err := users.NewUser(ids2.NewUserId(), "Test User", users.NewOwnership(ids2.NewAdminId()))
	require.NoError(t, err)
	user.PullDomainEvents() // clear creation event
	return user
}

// aUserWithId builds a User aggregate bound to the given userId.
func aUserWithId(t *testing.T, userId ids2.UserId) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, "Test User", users.NewOwnership(ids2.NewAdminId()))
	require.NoError(t, err)
	user.PullDomainEvents() // clear creation event
	return user
}
