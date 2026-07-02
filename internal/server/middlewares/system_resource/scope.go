package system_resource

import (
	"github.com/gin-gonic/gin"
	scopeCheckSystem "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/check_system_resource"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func Scope(qryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeId, ok := ids.ScopeIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_SCOPE_ID"))
			c.Abort()
			return
		}

		if _, err := qryBus.Dispatch(c.Request.Context(), scopeCheckSystem.NewQuery(scopeId)); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
