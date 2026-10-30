package path_params

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// UserIdFromAuthor seeds the UserId in context from the authenticated author.
// Used in self-service routes where the tenant_id is implicit (comes from the token).
func UserIdFromAuthor() gin.HandlerFunc {
	return func(c *gin.Context) {
		author, ok := vo.AuthorFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_AUTHOR"))
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), author.Id))
		c.Next()
	}
}
