package route_providers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	securityHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/security"
	userHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func registerUsers(g *gin.RouterGroup, p *Backoffice) {
	t := g.Group("/users")
	t.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.UserRole),
	)

	t.GET("",
		middlewares.ScopeMiddleware(vo.ScopeUsersRead),
		ownership.User(p.queryBus),
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, helpers.UserIdFromAuthor)),
	)
	t.PATCH("/name",
		middlewares.ScopeMiddleware(vo.ScopeUsersRename),
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, helpers.UserIdFromAuthor)),
	)
	t.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeUsersDelete),
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, helpers.UserIdFromAuthor)),
	)

	pr := t.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.PUT("",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesCreate),
		ownership.User(p.queryBus),
		p.executor.Run(middlewares.NewHandlerChain(p.pgxPool,
			func(c *gin.Context) error { c.Status(http.StatusCreated); return nil },
			prHandlers.NewCreateProtectedResourceHandler(
				p.commandBus,
				helpers.UserOwnershipFromAuthor,
				ids.ProtectedResourceIdFromContext,
				helpers.UserIdFromAuthor,
			),
			scopeHandlers.NewCreateScopesFromBodyHandler(
				p.commandBus,
				helpers.UserOwnershipFromAuthor,
				ids.ProtectedResourceIdFromContext,
				helpers.UserIdFromAuthor,
				helpers.GenerateScopeId,
			),
		)),
	)
	pr.GET("",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesRead),
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(p.queryBus, ids.ProtectedResourceIdFromContext)),
	)
	pr.PATCH("/name",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesRename),
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesDelete),
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
		middlewares.ScopeMiddleware(vo.ScopeScopesCreate),
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandBus,
			helpers.ProtectedResourceOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.ScopeIdFromContext,
		)),
	)
	scope.GET("",
		middlewares.ScopeMiddleware(vo.ScopeScopesRead),
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewGetScopeHandler(p.queryBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeScopesDelete),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.POST("/allowed-clients/:client_id",
		middlewares.ScopeMiddleware(vo.ScopeScopesAddAllowedClient),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("/allowed-clients/:client_id",
		middlewares.ScopeMiddleware(vo.ScopeScopesRemoveAllowedClient),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)

	client := t.Group("/clients/:client_id")
	client.Use(
		path_params.ClientIdParam(),
	)
	client.PUT("",
		middlewares.ScopeMiddleware(vo.ScopeClientsCreate),
		ownership.User(p.queryBus),
		p.executor.Run(middlewares.NewHandlerChain(p.pgxPool,
			func(c *gin.Context) error { c.JSON(http.StatusCreated, gin.H{}); return nil },
			clientHandlers.NewCreateClientHandler(
				p.commandBus, ids.ClientIdFromContext, helpers.UserIdFromAuthor, helpers.UserOwnershipFromAuthor,
			),
			securityHandlers.NewCreateClientAssertionKeyHandler(
				p.commandBus, ids.ClientIdFromContext,
			),
			securityHandlers.NewCreateClientSecretHandler(
				p.commandBus, ids.ClientIdFromContext,
			),
		)),
	)
	client.GET("",
		middlewares.ScopeMiddleware(vo.ScopeClientsRead),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	client.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeClientsDelete),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/flows",
		middlewares.ScopeMiddleware(vo.ScopeClientsAddFlow),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/redirect-uris",
		middlewares.ScopeMiddleware(vo.ScopeClientsAddRedirectUri),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, ids.ClientIdFromContext)),
	)
}
