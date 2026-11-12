package route_providers

import (
	"github.com/gin-gonic/gin"
	eventLogHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/event_logs"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

// registerEventLogs registers the admin-only read of the events persisted by
// the audit subscriber (events.audit_events).
func registerEventLogs(g *gin.RouterGroup, p *Backoffice) {
	e := g.Group("/events")
	e.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins, p.sessionRepository),
		middlewares.RoleMiddleware(vo.SystemAdministratorRole),
	)

	e.GET("",
		p.executor.Run(eventLogHandlers.NewListEventLogsHandler(p.queryBus)),
	)
}
