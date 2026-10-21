package change_name_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/change_name"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	adminId  ids2.UserId
	name     string
	callerId *ids2.UserId
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withAdminId(id ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.adminId = id }
}

func withCallerId(id ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.callerId = &id }
}

func aCommand(t *testing.T, opts ...commandOpt) *change_name.Command {
	t.Helper()
	b := &commandBuilder{
		adminId: ids2.NewUserId(),
		name:    "New Name",
	}
	for _, o := range opts {
		o(b)
	}
	callerId := b.adminId
	if b.callerId != nil {
		callerId = *b.callerId
	}
	return change_name.NewCommand(b.adminId, b.name, callerId)
}

func anAdmin(t *testing.T, adminId ids2.UserId, name string) *admins.Admin {
	t.Helper()
	admin, err := admins.NewAdmin(adminId, name, "", false)
	require.NoError(t, err)
	admin.PullDomainEvents()
	return admin
}
