package middlewares

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	"strings"

	"github.com/gin-gonic/gin"
	verifyJwt "github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
)

func AccessToken(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_TOKEN"))
			c.Abort()
			return
		}

		actor, err := query.DispatchTyped[vo.Actor](
			c.Request.Context(),
			queryBus,
			verifyJwt.NewQuery(token),
		)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(
			vo.WithActor(
				c.Request.Context(),
				actor,
			),
		)
		c.Next()
	}
}

func extractBearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
