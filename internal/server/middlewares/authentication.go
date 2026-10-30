package middlewares

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// Authentication picks how the access token is read from the request origin:
// requests from one of cookieOrigins (the backoffice frontend) are
// authenticated with the HttpOnly cookie; any other request must carry a
// Bearer token. Either way the token is then verified the same way.
func Authentication(queryBus query.Bus, expectedIssuer string, cookieOrigins []string, sessionRepositories ...auth.SessionRepository) gin.HandlerFunc {
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
					"sub":         nil,
					"sid":         nil,
					"system_role": nil,
					"memberships": nil,
					"picture":     nil,
				}),
		)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		subject, _ := claims["sub"].(string)
		if len(sessionRepositories) > 0 && sessionRepositories[0] != nil {
			sessionIDRaw, ok := claims["sid"].(string)
			if !ok || sessionIDRaw == "" {
				_ = c.Error(app_errors.UnauthorizedError("MISSING_SESSION_ID"))
				c.Abort()
				return
			}
			sessionID, err := ids.NewSessionIdFromString(sessionIDRaw)
			if err != nil {
				_ = c.Error(app_errors.UnauthorizedError("INVALID_SESSION_ID").WithCause(err, false))
				c.Abort()
				return
			}
			active, checkErr := sessionRepositories[0].IsActive(c.Request.Context(), sessionID)
			if checkErr != nil || !active {
				_ = c.Error(app_errors.UnauthorizedError("SESSION_REVOKED"))
				c.Abort()
				return
			}
		}
		systemRole, ok := claims["system_role"].(string)
		if !ok || systemRole == "" {
			_ = c.Error(app_errors.UnauthorizedError("MISSING_SYSTEM_ROLE"))
			c.Abort()
			return
		}
		memberships := make([]vo.SubscriptionMembership, 0)
		if values, ok := claims["memberships"].([]any); ok {
			for _, value := range values {
				membership, ok := value.(map[string]any)
				if !ok {
					_ = c.Error(app_errors.UnauthorizedError("INVALID_MEMBERSHIP"))
					c.Abort()
					return
				}
				rawID, idOK := membership["subscription_id"].(string)
				rawRole, roleOK := membership["subscription_role"].(string)
				if !idOK || !roleOK {
					_ = c.Error(app_errors.UnauthorizedError("INVALID_MEMBERSHIP"))
					c.Abort()
					return
				}
				id, idErr := ids.NewSubscriptionIdFromString(rawID)
				membershipRole, roleErr := vo.NewSubscriptionRole(rawRole)
				if idErr != nil || roleErr != nil {
					if idErr != nil {
						_ = c.Error(app_errors.UnauthorizedError("INVALID_MEMBERSHIP").WithCause(idErr, false))
					} else {
						_ = c.Error(app_errors.UnauthorizedError("INVALID_MEMBERSHIP").WithCause(roleErr, false))
					}
					c.Abort()
					return
				}
				memberships = append(memberships, vo.NewSubscriptionMembership(id, membershipRole))
			}
		}

		roleVO, err := vo.RoleFromSystemRoleClaim(systemRole)
		if err != nil {
			_ = c.Error(app_errors.UnauthorizedError("INVALID_SYSTEM_ROLE").WithCause(err, false))
			c.Abort()
			return
		}
		author, err := vo.NewUserAuthorWithMemberships(subject, subject, roleVO, memberships)
		if err != nil {
			_ = c.Error(app_errors.UnauthorizedError("INVALID_AUTHOR").WithCause(err, false))
			c.Abort()
			return
		}

		ctx := vo.WithAuthor(c.Request.Context(), author)
		if picture, ok := claims["picture"].(string); ok && picture != "" {
			ctx = vo.WithPicture(ctx, picture)
		}
		if rawScope, ok := claims["scope"].(string); ok {
			if scopes := vo.ParseScopes(rawScope); len(scopes) > 0 {
				ctx = vo.WithScopes(ctx, scopes)
			}
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
