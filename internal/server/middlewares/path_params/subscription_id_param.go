package path_params

import (
	"slices"

	"github.com/gin-gonic/gin"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	authors2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const subscriptionIdParam = "subscription_id"

func SubscriptionIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		subscriptionId, err := ids2.NewSubscriptionIdFromString(c.Param(subscriptionIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_SUBSCRIPTION_ID").WithDescription("invalid subscription id").WithCause(err, false))
			c.Abort()
			return
		}

		author, ok := vo2.AuthorFromContext[*authors2.UserAuthor](c.Request.Context())
		if !ok || !slices.Contains(author.SubscriptionIds, subscriptionId) {
			_ = c.Error(app_errors.ForbiddenError("FORBIDDEN_SUBSCRIPTION"))
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(
			ids2.WithSubscriptionId(c.Request.Context(), subscriptionId),
		)
		c.Next()
	}
}
