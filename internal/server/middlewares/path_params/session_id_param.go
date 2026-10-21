package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const sessionIdParam = "session_id"

func SessionIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := ids2.NewSessionIdFromString(c.Param(sessionIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_SESSION_ID").WithDescription("invalid session id").WithCause(err, false))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids2.WithSessionId(c.Request.Context(), sessionId),
		)
		c.Next()
	}
}
