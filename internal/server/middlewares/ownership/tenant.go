package ownership

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func User(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		author, ok := vo2.AuthorFromContext[*authors.AdminAuthor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_AUTHOR"))
			c.Abort()
			return
		}

		userId, ok := ids.UserIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_USER_ID"))
			c.Abort()
			return
		}

		result, err := query.DispatchTyped[users.Ownership](c.Request.Context(), queryBus,
			userCheckOwnership.NewQuery(userId.String(), author.Id),
		)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(vo2.WithOwnership(c.Request.Context(), result))
		c.Next()
	}
}
