package ownership

import (
	"github.com/gin-gonic/gin"
	clientCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func Client(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor, ok := vo2.ActorFromContext[*actors.AdminActor](c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_ACTOR"))
			c.Abort()
			return
		}

		clientId, ok := ids.ClientIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_CLIENT_ID"))
			c.Abort()
			return
		}

		userId, ok := ids.UserIdFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(app_errors.BadRequestError("MISSING_USER_ID"))
			c.Abort()
			return
		}

		result, err := queryBus.Dispatch(c.Request.Context(),
			clientCheckOwnership.NewQuery(
				clientId.String(),
				actor.Id,
				userId,
			))
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(vo2.WithOwnership(c.Request.Context(), result))
		c.Next()
	}
}
