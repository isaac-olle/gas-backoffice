package system_resource

import (
	"github.com/gin-gonic/gin"
	prCheckSystem "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/application/check_system_resource"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func ProtectedResource(qryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		prId, ok := ids.ProtectedResourceIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_PROTECTED_RESOURCE_ID"))
			c.Abort()
			return
		}

		if _, err := qryBus.Dispatch(c.Request.Context(), prCheckSystem.NewQuery(prId)); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
