package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

const scopeIdParam = "scope_id"

func ScopeIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeId, err := ids2.NewScopeIdFromString(c.Param(scopeIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_SCOPE_ID").WithCause(err))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			ids2.WithScopeId(c.Request.Context(), scopeId),
		)
		c.Next()
	}
}
