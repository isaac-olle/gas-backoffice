package check_tenant_access_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

func aTenant(t *testing.T, tenantId ids2.UserId, adminId ids2.UserId) *users.User {
	t.Helper()
	user, err := users.FromPrimitives(users.UserPrimitives{
		Id:      tenantId.String(),
		AdminId: adminId.String(),
		Name:    "Test Tenant",
		Role:    "tenant",
	})
	require.NoError(t, err)
	return user
}

func anAdmin(adminId ids2.UserId) *authors.UserAuthor {
	return &authors.UserAuthor{Id: adminId, Role: vo.AdminRole}
}

func aTenantAuthor(tenantId ids2.UserId) *authors.UserAuthor {
	return &authors.UserAuthor{Id: tenantId, Role: vo.TenantRole}
}
