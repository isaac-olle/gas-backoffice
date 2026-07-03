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
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, helpers.ClientIdFromActor)),
	)
	c.POST("/redirect-uris",
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, helpers.ClientIdFromActor)),
	)
	c.POST("/flows",
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, helpers.ClientIdFromActor)),
	)
}
