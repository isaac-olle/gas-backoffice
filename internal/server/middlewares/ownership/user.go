package ownership

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		author, ok := vo2.AuthorFromContext[*authors.UserAuthor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_USER_AUTHOR"))
			c.Abort()
			return
		}

		ownership := users.NewOwnership(author.Id)
		c.Request = c.Request.WithContext(vo2.WithOwnership(c.Request.Context(), ownership))
		c.Next()
	}
}
