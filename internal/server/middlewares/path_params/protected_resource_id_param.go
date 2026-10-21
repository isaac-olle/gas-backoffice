package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const protectedResourceIdParam = "protected_resource_id"

func ProtectedResourceIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		protectedResourceId, err := ids2.NewProtectedResourceIdFromString(c.Param(protectedResourceIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_PROTECTED_RESOURCE_ID").WithDescription("invalid protected resource id").WithCause(err, false))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids2.WithProtectedResourceId(c.Request.Context(), protectedResourceId),
		)
		c.Next()
	}
}
