package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

func RoleMiddleware(allowed ...vo2.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor, ok := vo2.ActorFromContext[vo2.Actor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_ACTOR"))
			c.Abort()
			return
		}

		for _, role := range allowed {
			if actor.EqualsRole(role) {
				c.Next()
				return
			}
		}

		_ = c.Error(app_errors.ForbiddenError("FORBIDDEN_ROLE"))
		c.Abort()
	}
}
