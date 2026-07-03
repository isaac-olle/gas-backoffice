package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	protocolHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol/authentication_methods"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	error_presenters "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_presenters"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type Protocol struct {
	executor   *middlewares.HandlerMiddleware
	commandBus command.Bus
	queryBus   query.Bus
}

func NewProtocol(hm *middlewares.HandlerMiddleware, cmdBus command.Bus, qryBus query.Bus) *Protocol {
	return &Protocol{executor: hm, commandBus: cmdBus, queryBus: qryBus}
}

func (p *Protocol) Register(group *gin.RouterGroup) {
	authorizePresenter := error_presenters.NewOpenIDAuthorizeErrorPresenter(
		session.ErrInvalidRedirectURI,
		session.ErrClientNotFound,
	)
	tokenPresenter := error_presenters.NewOAuth2TokenErrorPresenter()

	group.GET("/authorize",
		middlewares.HttpError(authorizePresenter),
		p.executor.Run(protocolHandler.NewAuthorizeHandler(p.commandBus, p.queryBus)),
	)
	group.POST("/token",
		middlewares.HttpError(tokenPresenter),
		p.executor.Run(protocolHandler.NewTokenAuthenticationHandler(
			authentication_methods.NewClientAuthMethodFactory(
				authentication_methods.NewBasicAuthMethod(p.queryBus),
				authentication_methods.NewJWTAssertionMethod(p.queryBus),
				authentication_methods.NewBodyCredentialsMethod(p.queryBus),
			),
		)),
		p.executor.Run(protocolHandler.NewTokenHandler(p.queryBus)),
	)
}
