package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	userHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// registerUsers registers self-service identity routes for all roles. The
// caller's identity comes from the access token — no path param needed.
// Resources (clients, protected resources, scopes) live under
// registerSubscriptionResources, scoped by :subscription_id.
func registerUsers(g *gin.RouterGroup, p *Backoffice) {
	u := g.Group("/users")
	u.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins),
		middlewares.RoleMiddleware(vo.AdminRole, vo.TenantRole),
		path_params.UserIdFromAuthor(),
	)

	u.GET("",
		middlewares.RoleMiddleware(vo.AdminRole),
		p.executor.Run(userHandlers.NewListUsersHandler(p.queryBus, ids.UserIdFromContext)),
	)
	u.GET("/me",
		p.executor.Run(userHandlers.NewGetCurrentUserHandler(p.queryBus, helpers.AuthorFromContext)),
	)
	u.DELETE("",
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandExecutor, ids.UserIdFromContext, helpers.AuthorFromContext)),
	)
	u.PATCH("/name",
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandExecutor, ids.UserIdFromContext, helpers.AuthorFromContext)),
	)
}
