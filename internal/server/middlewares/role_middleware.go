package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

func RoleMiddleware(allowed ...vo2.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		author, ok := vo2.AuthorFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_AUTHOR"))
			c.Abort()
			return
		}

		for _, role := range allowed {
			if author.EqualsRole(role) {
				c.Next()
				return
			}
		}

		_ = c.Error(app_errors.ForbiddenError("FORBIDDEN_ROLE"))
		c.Abort()
	}
}
