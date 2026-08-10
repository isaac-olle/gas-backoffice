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
	backofficePrId  string
	jwtIssuer       string
	publicBaseURL   string
}

func NewBackoffice(
	executor *middlewares.HandlerMiddleware,
	queryBus query.Bus,
	commandExecutor command.CommandExecutor,
	pgxPool *pgxpool.Pool,
	backofficePrId, jwtIssuer, publicBaseURL string,
) *Backoffice {
	return &Backoffice{
		executor:        executor,
		queryBus:        queryBus,
		commandExecutor: commandExecutor,
		pgxPool:         pgxPool,
		backofficePrId:  backofficePrId,
		jwtIssuer:       jwtIssuer,
		publicBaseURL:   publicBaseURL,
	}
}

func (p *Backoffice) Register(root *gin.RouterGroup) {
	v1 := root.Group("/v1")
	g := v1.Group("/backoffice")
	registerAdmins(g, p)
	registerUsers(g, p)

	// federated login (OIDC) — public, mounted at /v1/auth/... (not under /backoffice)
	registerAuth(v1, p)
}
