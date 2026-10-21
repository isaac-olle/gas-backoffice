package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Backoffice struct {
	executor        *middlewares.HandlerMiddleware
	queryBus        query.Bus
	commandExecutor command.CommandExecutor
	pgxPool         *pgxpool.Pool
	jwtIssuer       string
	cookieOrigins   []string
}

func NewBackoffice(
	executor *middlewares.HandlerMiddleware,
	queryBus query.Bus,
	commandExecutor command.CommandExecutor,
	pgxPool *pgxpool.Pool,
	jwtIssuer string,
	cookieOrigins []string,
) *Backoffice {
	return &Backoffice{
		executor:        executor,
		queryBus:        queryBus,
		commandExecutor: commandExecutor,
		pgxPool:         pgxPool,
		jwtIssuer:       jwtIssuer,
		cookieOrigins:   cookieOrigins,
	}
}

func (p *Backoffice) Register(root *gin.RouterGroup) {
	v1 := root.Group("/v1")
	g := v1.Group("/backoffice")
	registerAdmins(g, p)
	registerUsers(g, p)
	registerSubscriptions(g, p)
	registerSubscriptionResources(g, p)
	registerIdentityProviders(g, p)
	registerErrorLogs(g, p)
	registerAuth(g, p)
}
