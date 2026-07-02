package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/server/middlewares"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Backoffice struct {
	executor     *middlewares.HandlerMiddleware
	queryBus     query.Bus
	commandBus   command.Bus
	pgxPool      *pgxpool.Pool
	signingKeyId string
	jwtIssuer    string
}

func NewBackoffice(
	executor *middlewares.HandlerMiddleware,
	queryBus query.Bus,
	commandBus command.Bus,
	pgxPool *pgxpool.Pool,
	signingKeyId, jwtIssuer string,
) *Backoffice {
	return &Backoffice{
		executor:     executor,
		queryBus:     queryBus,
		commandBus:   commandBus,
		pgxPool:      pgxPool,
		signingKeyId: signingKeyId,
		jwtIssuer:    jwtIssuer,
	}
}

func (p *Backoffice) Register(group *gin.RouterGroup) {
	g := group.Group("/backoffice")
	registerAdmins(g, p)
	registerTenants(g, p)
	registerClients(g, p)
}
