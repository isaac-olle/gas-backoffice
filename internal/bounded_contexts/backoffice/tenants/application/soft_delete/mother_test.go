package soft_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/soft_delete"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/stretchr/testify/require"
)

type commandOpt func(*commandBuilder)

type commandBuilder struct {
	userId string
	author *authors.UserAuthor
}

func withUserId(userId string) commandOpt {
	return func(b *commandBuilder) { b.userId = userId }
}

func withAuthor(author *authors.UserAuthor) commandOpt {
	return func(b *commandBuilder) { b.author = author }
}

func aCommand(t *testing.T, opts ...commandOpt) *soft_delete.Command {
	t.Helper()
	b := &commandBuilder{userId: uuid.New().String()}
	for _, o := range opts {
		o(b)
	}
	id, err := ids2.NewUserIdFromString(b.userId)
	require.NoError(t, err)
	author := b.author
	if author == nil {
		author = &authors.UserAuthor{Id: id, Role: vo.TenantRole}
	}
	return soft_delete.NewCommand(id, author)
}

func aTenant(t *testing.T) *tenants.Tenant {
	t.Helper()
	userId := ids2.NewUserId()
	tenant, err := tenants.NewTenant(userId, "Test User", ids2.NewUserId(), "", false)
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}

func aTenantWithId(t *testing.T, userId ids2.UserId) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(userId, "Test User", ids2.NewUserId(), "", false)
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}
