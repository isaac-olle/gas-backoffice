package change_name_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
	"github.com/stretchr/testify/require"
)

// commandOpt is a functional option for commandBuilder.
type commandOpt func(*commandBuilder)

type commandBuilder struct {
	userId ids2.UserId
	name   string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withUserId(id ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = id }
}

// aCommand builds a change_name.Command using the Object Mother pattern.
func aCommand(t *testing.T, opts ...commandOpt) *change_name.Command {
	t.Helper()
	b := &commandBuilder{
		userId: ids2.NewUserId(),
		name:   "New Name",
	}
	for _, o := range opts {
		o(b)
	}
	return change_name.NewCommand(b.userId, b.name)
}

// aUser builds a User aggregate bound to the given userId and adminId.
func aUser(t *testing.T, userId ids2.UserId, adminId ids2.AdminId, name string) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, name, users.NewOwnership(adminId))
	require.NoError(t, err)
	user.PullDomainEvents() // clear creation event
	return user
}
