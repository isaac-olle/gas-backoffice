package change_name_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/change_name"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	userId ids2.UserId
	name   string
	author *authors.UserAuthor
}

func withName(name string) commandOpt {
	return func(b *commandBuilder) { b.name = name }
}

func withUserId(id ids2.UserId) commandOpt {
	return func(b *commandBuilder) { b.userId = id }
}

func withAuthor(author *authors.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
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
	author := b.author
	if author == nil {
		author = &authors.UserAuthor{Id: b.userId, Role: vo.TenantRole}
	}
	return change_name.NewCommand(b.userId, b.name, author)
}

func aTenant(t *testing.T, userId ids2.UserId, _ interface{}, name string) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(userId, name, ids2.NewUserId(), "", false)
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}
