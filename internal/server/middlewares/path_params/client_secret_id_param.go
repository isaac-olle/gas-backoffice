package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const clientSecretIdParam = "secret_id"

func ClientSecretIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretId, err := ids.NewClientSecretIdFromString(c.Param(clientSecretIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_CLIENT_SECRET_ID").WithDescription("invalid client secret id").WithCause(err))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids.WithClientSecretId(c.Request.Context(), secretId),
		)
		c.Next()
	}
}
