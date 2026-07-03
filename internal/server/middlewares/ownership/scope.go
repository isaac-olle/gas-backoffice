package ownership

import (
	"github.com/gin-gonic/gin"
	scopeCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func Scope(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor, ok := vo2.ActorFromContext[*actors.AdminActor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_ACTOR"))
			c.Abort()
			return
		}

		scopeId, ok := ids.ScopeIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_SCOPE_ID"))
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

		result, err := queryBus.Dispatch(c.Request.Context(), scopeCheckOwnership.NewQuery(
			scopeId.String(),
			protectedResourceId.String(),
			tenantId.String(),
			actor.Id.String(),
		))
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(vo2.WithOwnership(c.Request.Context(), result))
		c.Next()
	}
}
