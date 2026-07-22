package ownership

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	checkTenantAccess "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_tenant_access"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func Tenant(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, ok := ids.UserIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_TENANT_ID"))
			c.Abort()
			return
		}

		author, ok := vo2.AuthorFromContext[*authors.UserAuthor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_AUTHOR"))
			c.Abort()
			return
		}

		result, err := query.DispatchTyped[users.Ownership](c.Request.Context(), queryBus,
			checkTenantAccess.NewQuery(tenantId, author),
		)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(vo2.WithOwnership(c.Request.Context(), result))
		c.Next()
	}
}
