package create_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/create"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	id    ids2.UserId
	name  string
	email string
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withEmail(email string) commandOpt {
	return func(b *commandBuilder) { b.email = email }
}

func aCommand(t *testing.T, opts ...commandOpt) *create.Command {
	t.Helper()
	b := &commandBuilder{
		id:   ids2.NewUserId(),
		name: "Admin Name",
	}
	for _, o := range opts {
		o(b)
	}
	cmd := create.NewCommand(b.id, b.name, b.email)
	require.NotNil(t, cmd)
	return cmd
}
