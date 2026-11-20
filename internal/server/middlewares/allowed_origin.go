package middlewares

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func RequireAllowedOrigin(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := c.Request.Header.Values("Origin")
		if len(origins) != 1 || origins[0] == "" || origins[0] == "null" || !slices.Contains(allowedOrigins, origins[0]) {
			_ = c.Error(app_errors.ForbiddenError("ORIGIN_NOT_ALLOWED"))
			c.Abort()
			return
		}
		c.Next()
	}
}
