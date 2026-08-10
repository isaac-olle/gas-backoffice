package route_providers

import (
	"github.com/gin-gonic/gin"
	authHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
)

// registerAuth registers the federated-login (OIDC) routes.
//
// The whole OIDC dance lives on the backend (gas as relying party): the
// frontend only hits /login/:idp_id, the backend redirects the browser to the
// external IDP, and the IDP redirects back to /callback/:idp_id where the
// backend exchanges the code (with the client secret), resolves the user,
// creates the session and hands a token back to the frontend.
//
// PUBLIC (no AccessToken middleware): the browser reaches these unauthenticated.
// Mounted on the /v1 group directly (not under /backoffice) → /v1/auth/...
func registerAuth(g *gin.RouterGroup, p *Backoffice) {
	auth := g.Group("/auth")
	auth.GET("/login/:idp_id/:session_id",
		p.executor.Run(authHandlers.NewLoginHandler(p.commandExecutor, p.queryBus)),
	)
	auth.GET("/callback/:idp_id",
		p.executor.Run(authHandlers.NewOIDCCallbackHandler(p.commandExecutor, p.queryBus)),
	)
}
