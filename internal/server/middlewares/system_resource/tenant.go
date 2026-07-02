package system_resource

import (
	"github.com/gin-gonic/gin"
	tenantCheckSystem "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/check_system_resource"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func Tenant(qryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, ok := ids.TenantIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_TENANT_ID"))
			c.Abort()
			return
		}

		if _, err := qryBus.Dispatch(c.Request.Context(), tenantCheckSystem.NewQuery(tenantId)); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
