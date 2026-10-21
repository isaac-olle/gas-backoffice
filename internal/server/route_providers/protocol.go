package route_providers

import (
	"github.com/gin-gonic/gin"
	getJwks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/application/get_jwks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authorize"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"
	protocolHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	error_presenters "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_presenters"
)

type Protocol struct {
	executor             *middlewares.HandlerMiddleware
	authorizeService     *authorize.Service
	tokenService         *token.Service
	authClientService    *authenticate_client.Service
	jwksService          *getJwks.Service
	createSessionService *create_session.Service
	managers             []app_errors.ErrorManager
}

func NewProtocol(
	hm *middlewares.HandlerMiddleware,
	authorizeService *authorize.Service,
	tokenService *token.Service,
	authClientService *authenticate_client.Service,
	jwksService *getJwks.Service,
	createSessionService *create_session.Service,
	managers ...app_errors.ErrorManager,
) *Protocol {
	return &Protocol{
		executor:             hm,
		authorizeService:     authorizeService,
		tokenService:         tokenService,
		authClientService:    authClientService,
		jwksService:          jwksService,
		createSessionService: createSessionService,
		managers:             managers,
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
		p.executor.Run(protocolHandler.NewAuthorizeHandler(p.authorizeService)),
	)
	group.POST("/token",
		middlewares.HttpError(tokenPresenter, p.managers...),
		protocolHandler.TokenAuthentication(p.authClientService),
		p.executor.Run(protocolHandler.NewTokenHandler(p.tokenService)),
	)
	group.GET("/jwks",
		p.executor.Run(protocolHandler.NewJwksHandler(p.jwksService)),
	)

	// External IDP redirect target — :idp_id is the IDP reference name (e.g. google)
	root.GET("/v1/oidc/callback/:idp_id",
		p.executor.Run(protocolHandler.NewCallbackHandler(p.createSessionService)),
	)
}
