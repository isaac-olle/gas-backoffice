package route_providers

import (
	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
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
		middlewares.ScopeMiddleware("users:read"),
		ownership.User(p.queryBus),
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, helpers.UserIdFromAuthor)),
	)
	t.PATCH("/name",
		middlewares.ScopeMiddleware("users:rename"),
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, helpers.UserIdFromAuthor)),
	)
	t.DELETE("",
		middlewares.ScopeMiddleware("users:delete"),
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, helpers.UserIdFromAuthor)),
	)

	pr := t.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.PUT("",
		middlewares.ScopeMiddleware("protected_resources:create"),
		ownership.User(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(prHandlers.NewCreateProtectedResourceHandler(
			p.commandBus,
			helpers.UserOwnershipFromAuthor,
			ids.ProtectedResourceIdFromContext,
			helpers.UserIdFromAuthor,
		)),
		p.executor.Run(scopeHandlers.NewCreateScopesFromBodyHandler(
			p.commandBus,
			helpers.UserOwnershipFromAuthor,
			ids.ProtectedResourceIdFromContext,
			helpers.UserIdFromAuthor,
			helpers.GenerateScopeId,
		)),
	)
	pr.GET("",
		middlewares.ScopeMiddleware("protected_resources:read"),
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(p.queryBus, ids.ProtectedResourceIdFromContext)),
	)
	pr.PATCH("/name",
		middlewares.ScopeMiddleware("protected_resources:rename"),
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.DELETE("",
		middlewares.ScopeMiddleware("protected_resources:delete"),
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
		middlewares.ScopeMiddleware("scopes:create"),
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandBus,
			helpers.ProtectedResourceOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.ScopeIdFromContext,
		)),
	)
	scope.GET("",
		middlewares.ScopeMiddleware("scopes:read"),
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewGetScopeHandler(p.queryBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("",
		middlewares.ScopeMiddleware("scopes:delete"),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.POST("/allowed-clients/:client_id",
		middlewares.ScopeMiddleware("scopes:add_allowed_client"),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("/allowed-clients/:client_id",
		middlewares.ScopeMiddleware("scopes:remove_allowed_client"),
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
		middlewares.ScopeMiddleware("clients:create"),
		ownership.User(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, helpers.UserIdFromAuthor, helpers.UserOwnershipFromAuthor,
		)),
		p.executor.Run(clientHandlers.NewCreateClientAssertionKeyHandler(
			p.commandBus, ids.ClientIdFromContext,
		)),
		p.executor.Run(clientHandlers.NewCreateClientSecretHandler(
			p.commandBus, ids.ClientIdFromContext,
		)),
	)
	client.GET("",
		middlewares.ScopeMiddleware("clients:read"),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	client.DELETE("",
		middlewares.ScopeMiddleware("clients:delete"),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/flows",
		middlewares.ScopeMiddleware("clients:add_flow"),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/redirect-uris",
		middlewares.ScopeMiddleware("clients:add_redirect_uri"),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, ids.ClientIdFromContext)),
	)
}
