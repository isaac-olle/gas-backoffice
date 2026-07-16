package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const clientAssertionKeyIdParam = "assertion_key_id"

func ClientAssertionKeyIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyId, err := ids.NewClientAssertionKeyIdFromString(c.Param(clientAssertionKeyIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_ASSERTION_KEY_ID").WithDescription("invalid assertion key id").WithCause(err))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids.WithClientAssertionKeyId(c.Request.Context(), keyId),
		)
		c.Next()
	}
}
