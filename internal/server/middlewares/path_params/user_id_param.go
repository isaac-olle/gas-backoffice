package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const UserIdKey = "user_id"

func UserIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := ids2.NewUserIdFromString(c.Param(UserIdKey))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_USER_ID").WithDescription("invalid user id").WithCause(err, false))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids2.WithUserId(c.Request.Context(), userId),
		)
		c.Next()
	}
}
