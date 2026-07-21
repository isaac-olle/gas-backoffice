package change_name_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
	"github.com/stretchr/testify/require"
)

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

func aUser(t *testing.T, userId ids2.UserId, _ interface{}, name string) *users.User {
	t.Helper()
	user, err := users.NewUser(userId, name, users.NewOwnership(userId))
	require.NoError(t, err)
	user.PullDomainEvents()
	return user
}
