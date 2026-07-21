package ownership

import (
	"github.com/gin-gonic/gin"
	scopeCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func Scope(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		userId, ok := ids.UserIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_USER_ID"))
			c.Abort()
			return
		}

		result, err := queryBus.Dispatch(c.Request.Context(), scopeCheckOwnership.NewQuery(
			scopeId.String(),
			protectedResourceId.String(),
			userId.String(),
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
