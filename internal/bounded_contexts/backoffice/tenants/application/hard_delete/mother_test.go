package hard_delete_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/hard_delete"
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

// aCommand defaults to a tenant author that is self (author.Id == the
// command's own userId), so it succeeds the ownership check unless a test
// deliberately overrides withUserId/withAuthor to mismatch them.
func aCommand(t *testing.T, opts ...commandOpt) *hard_delete.Command {
	t.Helper()
	defaultUserId := ids2.NewUserId()
	b := &commandBuilder{
		userId: defaultUserId.String(),
		author: &authors.UserAuthor{Id: defaultUserId, Role: vo.TenantRole},
	}
	for _, o := range opts {
		o(b)
	}
	return hard_delete.NewCommand(b.userId, b.author)
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

// selfAuthor builds a tenant-role author that is the user itself (passes
// Tenant.IsSelf).
func selfAuthor(userId ids2.UserId) *authors.UserAuthor {
	return &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
}
