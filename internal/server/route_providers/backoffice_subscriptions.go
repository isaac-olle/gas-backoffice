package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	subscriptionHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

func registerSubscriptions(g *gin.RouterGroup, p *Backoffice) {
	r := g.Group("/subscriptions")
	r.Use(middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins, p.sessionRepository))
	h := subscriptionHandlers.NewHandler(p.commandExecutor, helpers.AuthorFromContext)
	// Kept as a server-side listing for consistency even though general users
	// already carry their subscription memberships in the JWT.
	r.GET("", p.executor.Run(subscriptionHandlers.NewListHandler(p.queryBus, helpers.AuthorFromContext)))
	r.Use(middlewares.RoleMiddleware(vo.SystemAdministratorRole))
	r.POST("", p.executor.Run(h))
	r.PUT("/:subscription_id", p.executor.Run(subscriptionHandlers.NewUpdateNameHandler(p.commandExecutor, helpers.AuthorFromContext)))
	r.POST("/:subscription_id/members/:user_id", p.executor.Run(h))
	r.DELETE("/:subscription_id/members/:user_id", p.executor.Run(h))
}
