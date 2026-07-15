package route_providers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/admins"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	securityHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/se
	userHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func registerAdmins(g *gin.RouterGroup, p *Backoffice) {
	a := g.Group("/admins")
	a.POST("/token",
		p.executor.Run(admins.NewTokenHandler(p.queryBus, p.backofficePrId, p.jwtIssuer)),
	)

	authenticated := a.Group("")
	authenticated.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.AdminRole),
	)
	authenticated.GET("",
		p.executor.Run(admins.NewGetAdminHandler(p.queryBus, helpers.AdminIdFromAuthor)),
	)

	user := authenticated.Group("/users/:user_id")
	user.Use(
		path_params.UserIdParam(),
	)
	user.PUT("",
		p.executor.Run(userHandlers.NewCreateUserHandler(
			p.commandBus, helpers.AdminIdFromAuthor, ids.UserIdFromContext,
		)),
	)
	user.DELETE("",
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, ids.UserIdFromContext)),
	)
	user.GET("",
		ownership.User(p.queryBus),
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, ids.UserIdFromContext)),
	)
	user.PATCH("/name",
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, ids.UserIdFromContext)),
	)

	pr := user.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.PUT("",
		ownership.User(p.queryBus),
		p.executor.Run(middlewares.NewHandlerChain(p.pgxPool,
			func(c *gin.Context) error { c.Status(http.StatusCreated); return nil },
			prHandlers.NewCreateProtectedResourceHandler(
				p.commandBus,
				helpers.UserOwnershipFromContext,
				ids.ProtectedResourceIdFromContext,
				ids.UserIdFromContext,
			),
			scopeHandlers.NewCreateScopesFromBodyHandler(
				p.commandBus,
				helpers.UserOwnershipFromContext,
				ids.ProtectedResourceIdFromContext,
				ids.UserIdFromContext,
				helpers.GenerateScopeId,
			),
		)),
	)
	pr.GET("",
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(
			p.queryBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.PATCH("/name",
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.DELETE("",
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewSoftDeleteProtectedResourceHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)

	scope := pr.Group("/scopes/:scope_id")
	scope.Use(
		path_params.ScopeIdParam(),
	)
	scope.PUT("",
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandBus,
			helpers.ProtectedResourceOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.ScopeIdFromContext,
		)),
	)
	scope.GET("",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewGetScopeHandler(p.queryBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("",
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.POST("/allowed-clients/:client_id",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("/allowed-clients/:client_id",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)

	client := user.Group("/clients/:client_id")
	client.Use(
		path_params.ClientIdParam(),
	)
	client.PUT("",
		ownership.User(p.queryBus),
		p.executor.Run(middlewares.NewHandlerChain(p.pgxPool,
			func(c *gin.Context) error { c.JSON(http.StatusCreated, gin.H{}); return nil },
			clientHandlers.NewCreateClientHandler(
				p.commandBus, ids.ClientIdFromContext, ids.UserIdFromContext, helpers.UserOwnershipFromContext,
			),
			securityHandlers.NewCreateClientAssertionKeyHandler(
				p.commandBus, ids.ClientIdFromContext,
			),
			securityHandlers.NewCreateClientSecretHandler(
				p.commandBus, ids.ClientIdFromContext,
			),
		)),
	)
	client.DELETE("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.GET("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	client.POST("/redirect-uris",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, ids.ClientIdFromContext)),
	)
}
