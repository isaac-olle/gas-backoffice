package route_providers

import (
	"github.com/gin-gonic/gin"
	errorLogHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

// registerErrorLogs registers the admin-only read of the errors persisted by
// the Postgres error manager (public.errors).
func registerErrorLogs(g *gin.RouterGroup, p *Backoffice) {
	e := g.Group("/errors")
	e.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins),
		middlewares.RoleMiddleware(vo.AdminRole),
	)

	e.GET("",
		p.executor.Run(errorLogHandlers.NewListErrorLogsHandler(p.queryBus)),
	)
}
