package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

const TenantIdKey = "tenant_id"

func TenantIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, err := ids2.NewTenantIdFromString(c.Param(TenantIdKey))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_TENANT_ID").WithCause(err))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids2.WithTenantId(c.Request.Context(), tenantId),
		)
		c.Next()
	}
}
