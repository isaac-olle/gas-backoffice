package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Backoffice struct {
	executor       *middlewares.HandlerMiddleware
	queryBus       query.Bus
	commandBus     command.Bus
	pgxPool        *pgxpool.Pool
	backofficePrId string
	jwtIssuer      string
}

func NewBackoffice(
	executor *middlewares.HandlerMiddleware,
	queryBus query.Bus,
	commandBus command.Bus,
	pgxPool *pgxpool.Pool,
	backofficePrId, jwtIssuer string,
) *Backoffice {
	return &Backoffice{
		executor:       executor,
		queryBus:       queryBus,
		commandBus:     commandBus,
		pgxPool:        pgxPool,
		backofficePrId: backofficePrId,
		jwtIssuer:      jwtIssuer,
	}
}

func (p *Backoffice) Register(group *gin.RouterGroup) {
	g := group.Group("/backoffice")
	registerUsers(g, p)
	registerClients(g, p)
}
