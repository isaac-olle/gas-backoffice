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
	r.Use(middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins), middlewares.RoleMiddleware(vo.AdminRole))
	h := subscriptionHandlers.NewHandler(p.commandExecutor, helpers.AuthorFromContext)
	r.POST("", p.executor.Run(h))
	r.POST("/:subscription_id/members/:user_id", p.executor.Run(h))
	r.DELETE("/:subscription_id/members/:user_id", p.executor.Run(h))
}
