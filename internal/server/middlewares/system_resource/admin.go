package system_resource

import (
	"github.com/gin-gonic/gin"
	adminCheckSystem "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/check_system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func Admin(qryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminId, ok := ids.AdminIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_ADMIN_ID"))
			c.Abort()
			return
		}

		if _, err := qryBus.Dispatch(c.Request.Context(), adminCheckSystem.NewQuery(adminId)); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
