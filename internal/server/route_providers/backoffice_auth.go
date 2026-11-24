package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/get_csrf_token"
	authInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	authHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

// registerLogin handles top-level browser navigation to start federated login.
// Mounted under /v1/backoffice → /v1/backoffice/auth/...
func registerLogin(g *gin.RouterGroup, p *Backoffice) {
	g.GET("/auth/login",
		p.executor.Run(authHandlers.NewLoginHandler(p.queryBus)),
	)
}

func registerAuth(g *gin.RouterGroup, p *Backoffice) {
	auth := g.Group("/auth", middlewares.RequireAllowedOrigin(p.cookieOrigins))
	auth.POST("/callback",
		p.executor.Run(authHandlers.NewCallbackHandler(p.queryBus)),
	)
	csrfTokens := get_csrf_token.NewService(p.sessionRepository, authInfra.RandomCSRFTokenGenerator{})
	refreshVerifier := authInfra.NewRefreshTokenVerifier(p.queryBus, p.jwtIssuer)
	auth.GET("/csrf",
		p.executor.Run(authHandlers.NewCSRFHandler(csrfTokens, refreshVerifier, p.queryBus, p.jwtIssuer)),
	)
	auth.POST("/refresh",
		p.executor.Run(authHandlers.NewRefreshHandler(p.queryBus)),
	)
	auth.POST("/logout",
		p.executor.Run(authHandlers.NewLogoutHandler(p.commandExecutor, p.queryBus, p.jwtIssuer, p.cookieAuthenticator)),
	)
}
