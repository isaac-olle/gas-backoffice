package route_providers

import (
	"github.com/gin-gonic/gin"
	authHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
)

// registerAuth registers the federated-login (OIDC) routes: the IDP redirects
// the browser to the frontend with ?code&state, the frontend POSTs them here,
// and the backend exchanges the code (with the client secret), resolves the
// user, creates the session and sets it as an HttpOnly cookie.
//
// PUBLIC (no AccessToken middleware): the browser reaches this unauthenticated.
// Mounted under /v1/backoffice → /v1/backoffice/auth/...
func registerAuth(g *gin.RouterGroup, p *Backoffice) {
	auth := g.Group("/auth")
	auth.GET("/login",
		p.executor.Run(authHandlers.NewLoginHandler(p.queryBus)),
	)
	auth.POST("/callback",
		p.executor.Run(authHandlers.NewCallbackHandler(p.queryBus)),
	)
}
