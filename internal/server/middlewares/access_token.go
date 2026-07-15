package middlewares

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"strings"

	"github.com/gin-gonic/gin"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
)

func AccessToken(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_TOKEN"))
			c.Abort()
			return
		}

		result, err := query.DispatchTyped[*verifyJwt.VerifyJWTResult](
			c.Request.Context(),
			queryBus,
			verifyJwt.NewQuery(token),
		)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		ctx := vo.WithAuthor(c.Request.Context(), result.Author)
		if len(result.Scopes) > 0 {
			ctx = vo.WithScopes(ctx, result.Scopes)
		}
		c.Request = c.Request.WithContext(ctx)
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
