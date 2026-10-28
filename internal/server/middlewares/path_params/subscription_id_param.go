package path_params

import (
	"github.com/gin-gonic/gin"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const subscriptionIdParam = "subscription_id"

// SubscriptionIdParam parses :subscription_id and stores it in the request
// context. It only validates the path value; it does not authorize access.
// Authorization is deliberately handled by middlewares.SubscriptionMiddleware.
func SubscriptionIdParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		subscriptionId, err := ids2.NewSubscriptionIdFromString(c.Param(subscriptionIdParam))
		if err != nil {
			_ = c.Error(app_errors.BadRequestError("INVALID_SUBSCRIPTION_ID").WithDescription("invalid subscription id").WithCause(err, false))
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(
			ids2.WithSubscriptionId(c.Request.Context(), subscriptionId),
		)
		c.Next()
	}
}
