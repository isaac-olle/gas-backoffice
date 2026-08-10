package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	protocolHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol/authentication_methods"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	error_presenters "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_presenters"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type Protocol struct {
	executor        *middlewares.HandlerMiddleware
	commandExecutor command.CommandExecutor
	queryBus        query.Bus
	managers        []app_errors.ErrorManager
}

func NewProtocol(hm *middlewares.HandlerMiddleware, commandExecutor command.CommandExecutor, qryBus query.Bus, managers ...app_errors.ErrorManager) *Protocol {
	return &Protocol{
		executor:        hm,
		commandExecutor: commandExecutor,
		queryBus:        qryBus,
		managers:        managers,
	}
}

func (p *Protocol) Register(root *gin.RouterGroup) {
	authorizePresenter := error_presenters.NewOpenIDAuthorizeErrorPresenter(
		session.ErrInvalidRedirectURI,
		session.ErrClientNotFound,
	)
	tokenPresenter := error_presenters.NewOAuth2TokenErrorPresenter()

	// OIDC provider surface — unversioned protocol endpoints at /oidc/*
	group := root.Group("/oidc")

	group.GET("/authorize",
		middlewares.HttpError(authorizePresenter, p.managers...),
		p.executor.Run(protocolHandler.NewAuthorizeHandler(p.commandExecutor, p.queryBus)),
	)
	group.POST("/token",
		middlewares.HttpError(tokenPresenter, p.managers...),
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
