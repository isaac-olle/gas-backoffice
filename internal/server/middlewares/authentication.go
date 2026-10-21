package middlewares

import (
	"slices"

	"github.com/gin-gonic/gin"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

// Authentication picks how the access token is read from the request origin:
// requests from one of cookieOrigins (the backoffice frontend) are
// authenticated with the HttpOnly cookie; any other request must carry a
// Bearer token. Either way the token is then verified the same way.
func Authentication(queryBus query.Bus, expectedIssuer string, cookieOrigins []string) gin.HandlerFunc {
	cookie := authenticators.NewCookieAuthenticator()
	bearer := authenticators.NewBearerAuthenticator()

	return func(c *gin.Context) {
		var authenticator authenticators.Authenticator = bearer
		if slices.Contains(cookieOrigins, c.GetHeader("Origin")) {
			authenticator = cookie
		}

		token, err := authenticator.Token(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		claims, err := query.DispatchTyped[map[string]any](
			c.Request.Context(),
			queryBus,
			verifyJwt.NewQuery(token).
				WithExpectedClaim("iss", expectedIssuer).
				WithExpectedClaim("aud", expectedIssuer).
				WithVerifyExpiration().
				WithClaimsToReturn(map[string]any{
					"sub":              nil,
					"role":             nil,
					"admin_id":         nil,
					"subscription_ids": nil,
				}),
		)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		subject, _ := claims["sub"].(string)
		roleStr, _ := claims["role"].(string)
		adminIdStr, _ := claims["admin_id"].(string)
		subscriptionIdStrings := make([]string, 0)
		if values, ok := claims["subscription_ids"].([]any); ok {
			for _, value := range values {
				if id, ok := value.(string); ok {
					subscriptionIdStrings = append(subscriptionIdStrings, id)
				}
			}
		}

		roleVO, err := vo.NewRole(roleStr)
		if err != nil {
			_ = c.Error(app_errors.UnauthorizedError("INVALID_ROLE").WithCause(err, false))
			c.Abort()
			return
		}

		author, err := authors.NewAuthorForRole(roleVO, subject, adminIdStr, subscriptionIdStrings)
		if err != nil {
			_ = c.Error(app_errors.UnauthorizedError("INVALID_AUTHOR").WithCause(err, false))
			c.Abort()
			return
		}

		ctx := vo.WithAuthor(c.Request.Context(), author)
		if rawScope, ok := claims["scope"].(string); ok {
			if scopes := vo.ParseScopes(rawScope); len(scopes) > 0 {
				ctx = vo.WithScopes(ctx, scopes)
			}
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
