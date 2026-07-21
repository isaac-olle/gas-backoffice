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
	u := g.Group("/users")

	// unauthenticated
	u.POST("/token",
		p.executor.Run(userHandlers.NewTokenHandler(p.queryBus, p.backofficePrId, p.jwtIssuer)),
	)

	// ── admin management routes (role=admin) ─────────────────────────────────
	admin := u.Group("/managed/:user_id")
	admin.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.AdminRole),
		path_params.UserIdParam(),
	)
	admin.PUT("",
		p.executor.Run(userHandlers.NewCreateUserHandler(
			p.commandBus, ids.UserIdFromContext,
		)),
	)
	admin.GET("",
		ownership.User(p.queryBus),
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, ids.UserIdFromContext)),
	)
	admin.DELETE("",
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, ids.UserIdFromContext)),
	)
	admin.PATCH("/name",
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, ids.UserIdFromContext)),
	)

	adminPR := admin.Group("/protected-resources/:protected_resource_id")
	adminPR.Use(path_params.ProtectedResourceIdParam())
	adminPR.PUT("",
		ownership.User(p.queryBus),
		p.executor.Run(middlewares.NewHandlerChain(p.pgxPool,
			func(c *gin.Context) error { c.Status(http.StatusCreated); return nil },
			prHandlers.NewCreateProtectedResourceHandler(
				p.commandBus,
				ids.ProtectedResourceIdFromContext,
				ids.UserIdFromContext,
			),
			scopeHandlers.NewCreateScopesFromBodyHandler(
				p.commandBus,
				ids.ProtectedResourceIdFromContext,
				ids.UserIdFromContext,
				helpers.GenerateScopeId,
			),
		)),
	)
	adminPR.GET("",
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(
			p.queryBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	adminPR.PATCH("/name",
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	adminPR.DELETE("",
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewSoftDeleteProtectedResourceHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)

	adminScope := adminPR.Group("/scopes/:scope_id")
	adminScope.Use(path_params.ScopeIdParam())
	adminScope.PUT("",
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandBus,
			helpers.ProtectedResourceOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.ScopeIdFromContext,
		)),
	)
	adminScope.GET("",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewGetScopeHandler(p.queryBus, ids.ScopeIdFromContext)),
	)
	adminScope.DELETE("",
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	adminScope.POST("/allowed-clients/:client_id",
		path_params.ClientIdParam(),
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext)),
	)
	adminScope.DELETE("/allowed-clients/:client_id",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)

	adminClient := admin.Group("/clients/:client_id")
	adminClient.Use(path_params.ClientIdParam())
	adminClient.PUT("",
		ownership.User(p.queryBus),
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, ids.UserIdFromContext,
		)),
	)
	adminClient.PUT("/assertion-keys/:assertion_key_id",
		path_params.ClientAssertionKeyIdParam(),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientAssertionKeyHandler(
			p.commandBus, ids.ClientIdFromContext, ids.ClientAssertionKeyIdFromContext,
		)),
	)
	adminClient.PUT("/secrets/:secret_id",
		path_params.ClientSecretIdParam(),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientSecretHandler(
			p.commandBus, ids.ClientIdFromContext, ids.ClientSecretIdFromContext,
		)),
	)
	adminClient.DELETE("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	adminClient.POST("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	adminClient.DELETE("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewRemoveAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	adminClient.GET("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	adminClient.POST("/redirect-uris",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, ids.ClientIdFromContext)),
	)

	// ── user self-service routes (role=user) ──────────────────────────────────
	self := u.Group("")
	self.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.UserRole),
	)
	self.GET("",
		middlewares.ScopeMiddleware(vo.ScopeUsersRead),
		ownership.User(p.queryBus),
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, helpers.UserIdFromAuthor)),
	)
	self.PATCH("/name",
		middlewares.ScopeMiddleware(vo.ScopeUsersRename),
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, helpers.UserIdFromAuthor)),
	)
	self.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeUsersDelete),
		ownership.User(p.queryBus),
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, helpers.UserIdFromAuthor)),
	)

	selfPR := self.Group("/protected-resources/:protected_resource_id")
	selfPR.Use(path_params.ProtectedResourceIdParam())
	selfPR.PUT("",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesCreate),
		ownership.User(p.queryBus),
		p.executor.Run(middlewares.NewHandlerChain(p.pgxPool,
			func(c *gin.Context) error { c.Status(http.StatusCreated); return nil },
			prHandlers.NewCreateProtectedResourceHandler(
				p.commandBus,
				ids.ProtectedResourceIdFromContext,
				helpers.UserIdFromAuthor,
			),
			scopeHandlers.NewCreateScopesFromBodyHandler(
				p.commandBus,
				ids.ProtectedResourceIdFromContext,
				helpers.UserIdFromAuthor,
				helpers.GenerateScopeId,
			),
		)),
	)
	selfPR.GET("",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesRead),
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(p.queryBus, ids.ProtectedResourceIdFromContext)),
	)
	selfPR.PATCH("/name",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesRename),
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	selfPR.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeProtectedResourcesDelete),
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewSoftDeleteProtectedResourceHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)

	selfScope := selfPR.Group("/scopes/:scope_id")
	selfScope.Use(path_params.ScopeIdParam())
	selfScope.PUT("",
		middlewares.ScopeMiddleware(vo.ScopeScopesCreate),
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandBus,
			helpers.ProtectedResourceOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.ScopeIdFromContext,
		)),
	)
	selfScope.GET("",
		middlewares.ScopeMiddleware(vo.ScopeScopesRead),
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewGetScopeHandler(p.queryBus, ids.ScopeIdFromContext)),
	)
	selfScope.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeScopesDelete),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(p.commandBus, ids.ScopeIdFromContext)),
	)
	selfScope.POST("/allowed-clients/:client_id",
		path_params.ClientIdParam(),
		middlewares.ScopeMiddleware(vo.ScopeScopesAddAllowedClient),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext)),
	)
	selfScope.DELETE("/allowed-clients/:client_id",
		middlewares.ScopeMiddleware(vo.ScopeScopesRemoveAllowedClient),
		ownership.Scope(p.queryBus),
		system_resource.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)

	selfClient := self.Group("/clients/:client_id")
	selfClient.Use(path_params.ClientIdParam())
	selfClient.PUT("",
		middlewares.ScopeMiddleware(vo.ScopeClientsCreate),
		ownership.User(p.queryBus),
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, helpers.UserIdFromAuthor,
		)),
	)
	selfClient.PUT("/assertion-keys/:assertion_key_id",
		path_params.ClientAssertionKeyIdParam(),
		middlewares.ScopeMiddleware(vo.ScopeClientsCreate),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientAssertionKeyHandler(
			p.commandBus, ids.ClientIdFromContext, ids.ClientAssertionKeyIdFromContext,
		)),
	)
	selfClient.PUT("/secrets/:secret_id",
		path_params.ClientSecretIdParam(),
		middlewares.ScopeMiddleware(vo.ScopeClientsCreate),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientSecretHandler(
			p.commandBus, ids.ClientIdFromContext, ids.ClientSecretIdFromContext,
		)),
	)
	selfClient.GET("",
		middlewares.ScopeMiddleware(vo.ScopeClientsRead),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	selfClient.DELETE("",
		middlewares.ScopeMiddleware(vo.ScopeClientsDelete),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	selfClient.POST("/flows",
		middlewares.ScopeMiddleware(vo.ScopeClientsAddFlow),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	selfClient.DELETE("/flows",
		middlewares.ScopeMiddleware(vo.ScopeClientsRemoveFlow),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewRemoveAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
	)
	selfClient.POST("/redirect-uris",
		middlewares.ScopeMiddleware(vo.ScopeClientsAddRedirectUri),
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandBus, ids.ClientIdFromContext)),
	)
}
