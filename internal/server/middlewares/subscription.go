package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/subscription"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// SubscriptionMiddleware resolves the subscription through the injected
// resolver and authorizes the caller against the memberships in the verified
// access token. It does not know whether the subscription came from a path or
// from an existing aggregate.
//
// Resource-specific loading is delegated to the resolver supplied by the
// caller; the middleware itself only authorizes the resolved subscription.
func SubscriptionMiddleware(resolve subscription.Resolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		subscriptionId, err := resolve.Resolve(c.Request.Context(), c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		author, ok := vo.AuthorFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.ForbiddenError("FORBIDDEN_SUBSCRIPTION"))
			c.Abort()
			return
		}

		if author.Role.IsSystemAdministrator() || hasSubscriptionMembership(author, subscriptionId) {
			c.Request = c.Request.WithContext(ids.WithSubscriptionId(c.Request.Context(), subscriptionId))
			c.Next()
			return
		}

		_ = c.Error(app_errors.ForbiddenError("FORBIDDEN_SUBSCRIPTION"))
		c.Abort()
	}
}

func hasSubscriptionMembership(author *vo.UserAuthor, subscriptionId ids.SubscriptionId) bool {
	_, ok := author.MembershipFor(subscriptionId)
	return ok
}
