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
		ownership.User(p.queryBus),
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, helpers.UserIdFromActor)),
	)
	t.PATCH("/name",
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, helpers.UserIdFromActor)),
	)
	t.DELETE("",
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, helpers.UserIdFromActor)),
	)

	pr := t.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.PUT("",
		ownership.User(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(prHandlers.NewCreateProtectedResourceHandler(
			p.commandBus,
			helpers.UserOwnershipFromActor,
			ids.ProtectedResourceIdFromContext,
			helpers.UserIdFromActor,
		)),
		p.executor.Run(scopeHandlers.NewCreateScopesFromBodyHandler(
			p.commandBus,
			helpers.UserOwnershipFromActor,
			ids.ProtectedResourceIdFromContext,
			helpers.UserIdFromActor,
			helpers.GenerateScopeId,
		)),
	)
	pr.GET("",
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(p.queryBus, ids.ProtectedResourceIdFromContext)),
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
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	scope.DELETE("/allowed-clients/:client_id",
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
		ownership.User(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, helpers.UserIdFromActor, helpers.UserOwnershipFromActor,
		)),
		p.executor.Run(clientHandlers.NewCreateClientAssertionKeyHandler(
			p.commandBus, ids.ClientIdFromContext,
		)),
		p.executor.Run(clientHandlers.NewCreateClientSecretHandler(
			p.commandBus, ids.ClientIdFromContext,
		)),
	)
	client.GET("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	client.DELETE("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	client.POST("/redirect-uris",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, ids.ClientIdFromContext)),
	)
}
