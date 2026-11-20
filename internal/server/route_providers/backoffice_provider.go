package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/verify_csrf_token"
	authInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Backoffice struct {
	executor            *middlewares.HandlerMiddleware
	queryBus            query.Bus
	commandExecutor     command.CommandExecutor
	pgxPool             *pgxpool.Pool
	jwtIssuer           string
	cookieOrigins       []string
	cookieAuthenticator *authenticators.CookieAuthenticator
	sessionRepository   *authInfra.PostgresSessionRepository
}

func NewBackoffice(
	executor *middlewares.HandlerMiddleware,
	queryBus query.Bus,
	commandExecutor command.CommandExecutor,
	pgxPool *pgxpool.Pool,
	jwtIssuer string,
	cookieOrigins []string,
	encrypter persistance.Encrypter,
	requireCSRF bool,
) *Backoffice {
	sessionRepository := authInfra.NewPostgresSessionRepository(pgxPool, encrypter)
	var csrfVerifier auth.CSRFTokenVerifier
	if requireCSRF {
		csrfVerifier = verify_csrf_token.NewService(sessionRepository)
	}
	return &Backoffice{
		executor:            executor,
		queryBus:            queryBus,
		commandExecutor:     commandExecutor,
		pgxPool:             pgxPool,
		jwtIssuer:           jwtIssuer,
		cookieOrigins:       cookieOrigins,
		sessionRepository:   sessionRepository,
		cookieAuthenticator: authenticators.NewCookieAuthenticator(csrfVerifier),
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
	registerEventLogs(g, p)
	registerAuth(g, p)
}
