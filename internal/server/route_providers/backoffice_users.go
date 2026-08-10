package route_providers

import (
	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	securityHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/security"
	sessionHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/sessions"
	userHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// registerUsers registers self-service routes for all roles.
// The caller's identity comes from the access token — no path param needed.
func registerUsers(g *gin.RouterGroup, p *Backoffice) {
	u := g.Group("/users")
	u.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.AdminRole, vo.TenantRole),
		path_params.UserIdFromAuthor(),
	)

	u.GET("",
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, ids.UserIdFromContext)),
	)
	u.DELETE("",
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandExecutor, ids.UserIdFromContext)),
	)
	u.PATCH("/name",
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandExecutor, ids.UserIdFromContext)),
	)

	session := u.Group("/sessions/:session_id")
	session.Use(path_params.SessionIdParam())
	session.PUT("",
		p.executor.Run(sessionHandlers.NewCreateSessionHandler(
			p.commandExecutor, ids.SessionIdFromContext, ids.UserIdFromContext,
		)),
	)
	session.GET("",
		p.executor.Run(sessionHandlers.NewGetSessionHandler(p.queryBus, ids.SessionIdFromContext)),
	)

	pr := u.Group("/protected-resources/:protected_resource_id")
	pr.Use(path_params.ProtectedResourceIdParam())
	pr.PUT("",
		p.executor.Run(
			prHandlers.NewCreateProtectedResourceHandler(
				p.commandExecutor,
				ids.ProtectedResourceIdFromContext,
				ids.UserIdFromContext,
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
			p.commandExecutor, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.DELETE("",
		ownership.ProtectedResource(p.queryBus),
		system_resource.ProtectedResource(p.queryBus),
		p.executor.Run(prHandlers.NewSoftDeleteProtectedResourceHandler(
			p.commandExecutor, ids.ProtectedResourceIdFromContext,
		)),
	)

	scope := pr.Group("/scopes/:scope_id")
	scope.Use(path_params.ScopeIdParam())
	scope.PUT("",
		ownership.ProtectedResource(p.queryBus),
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandExecutor,
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
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(p.commandExecutor, ids.ScopeIdFromContext)),
	)
	scope.POST("/allowed-clients/:client_id",
		path_params.ClientIdParam(),
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)
	scope.DELETE("/allowed-clients/:client_id",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)

	client := u.Group("/clients/:client_id")
	client.Use(path_params.ClientIdParam())
	client.PUT("",
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.UserIdFromContext,
		)),
	)
	client.PUT("/assertion-keys/:assertion_key_id",
		path_params.ClientAssertionKeyIdParam(),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientAssertionKeyHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.ClientAssertionKeyIdFromContext,
		)),
	)
	client.PUT("/secrets/:secret_id",
		path_params.ClientSecretIdParam(),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientSecretHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.ClientSecretIdFromContext,
		)),
	)
	client.DELETE("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(p.commandExecutor, ids.ClientIdFromContext)),
	)
	client.POST("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddAllowedFlowHandler(p.commandExecutor, ids.ClientIdFromContext)),
	)
	client.DELETE("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewRemoveAllowedFlowHandler(p.commandExecutor, ids.ClientIdFromContext)),
	)
	client.GET("",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewGetClientHandler(p.queryBus, ids.ClientIdFromContext)),
	)
	client.POST("/redirect-uris",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(p.commandExecutor, ids.ClientIdFromContext)),
	)
}
