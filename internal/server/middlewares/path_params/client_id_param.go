package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const clientIdParam = "client_id"

func ClientIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId, err := ids2.NewClientIdFromString(c.Param(clientIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_CLIENT_ID").WithDescription("invalid client id").WithCause(err, false))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids2.WithClientId(c.Request.Context(), clientId),
		)
		c.Next()
	}
}
