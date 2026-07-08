package system_resource

import (
	"github.com/gin-gonic/gin"
	userCheckSystem "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func User(qryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := ids.UserIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_USER_ID"))
			c.Abort()
			return
		}

		if _, err := qryBus.Dispatch(c.Request.Context(), userCheckSystem.NewQuery(userId)); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
