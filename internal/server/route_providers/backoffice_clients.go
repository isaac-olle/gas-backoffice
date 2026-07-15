package route_providers

import (
	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

func registerClients(g *gin.RouterGroup, p *Backoffice) {
	c := g.Group("/clients")
	c.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.ClientRole),
	)

	c.GET("",
		middlewares.ScopeMiddleware("clients:read"),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, helpers.ClientIdFromAuthor)),
	)
	c.POST("/redirect-uris",
		middlewares.ScopeMiddleware("clients:add_redirect_uri"),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, helpers.ClientIdFromAuthor)),
	)
	c.POST("/flows",
		middlewares.ScopeMiddleware("clients:add_flow"),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, helpers.ClientIdFromAuthor)),
	)
}
