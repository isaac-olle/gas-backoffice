package route_providers

import (
	"github.com/gin-gonic/gin"
	idpHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

// registerIdentityProviders registers admin-only identity provider management routes.
func registerIdentityProviders(g *gin.RouterGroup, p *Backoffice) {
	idp := g.Group("/identity-providers")
	idp.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins),
		middlewares.RoleMiddleware(vo.AdminRole),
	)

	idp.POST("",
		p.executor.Run(idpHandlers.NewCreateIdentityProviderHandler(p.commandExecutor)),
	)
}
