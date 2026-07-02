package ownership

import (
	"github.com/gin-gonic/gin"
	prCheckOwnership "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	vo2 "github.com/isaac-olle/gas/internal/shared/domain/vo"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/actors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func ProtectedResource(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor, ok := vo2.ActorFromContext[*actors.AdminActor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_ACTOR"))
			c.Abort()
			return
		}

		protectedResourceId, ok := ids.ProtectedResourceIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_PROTECTED_RESOURCE_ID"))
			c.Abort()
			return
		}

		tenantId, ok := ids.TenantIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_TENANT_ID"))
			c.Abort()
			return
		}

		result, err := queryBus.Dispatch(
			c.Request.Context(),
			prCheckOwnership.NewQuery(protectedResourceId.String(), tenantId.String(), actor.Id),
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
