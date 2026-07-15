package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

func ScopeMiddleware(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !vo.HasScope(c.Request.Context(), scope) {
			_ = c.Error(app_errors.ForbiddenError("INSUFFICIENT_SCOPE").
				WithDescription("token does not have required scope: " + scope))
			c.Abort()
			return
		}
		c.Next()
	}
}
