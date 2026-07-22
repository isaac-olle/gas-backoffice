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

// registerUsers registers self-service routes for all roles.
// The caller's identity comes from the access token — no path param needed.
func registerUsers(g *gin.RouterGroup, p *Backoffice) {
	u := g.Group("/users")
	u.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.AdminRole, vo.TenantRole),
		path_params.UserIdFromAuthor(),
	)

	u.POST("/token",
		p.executor.Run(userHandlers.NewTokenHandler(p.queryBus, p.backofficePrId, p.jwtIssuer)),
	)
	u.GET("",
		p.executor.Run(userHandlers.NewGetUserHandler(p.queryBus, ids.UserIdFromContext)),
	)
	u.DELETE("",
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewSoftDeleteUserHandler(p.commandBus, ids.UserIdFromContext)),
	)
	u.PATCH("/name",
		system_resource.User(p.queryBus),
		p.executor.Run(userHandlers.NewChangeUserNameHandler(p.commandBus, ids.UserIdFromContext)),
	)

	pr := u.Group("/protected-resources/:protected_resource_id")
	pr.Use(path_params.ProtectedResourceIdParam())
	pr.PUT("",
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
	scope.Use(path_params.ScopeIdParam())
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
		path_params.ClientIdParam(),
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(
			p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)
	scope.DELETE("/allowed-clients/:client_id",
		ownership.Scope(p.queryBus),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)),
	)

	client := u.Group("/clients/:client_id")
	client.Use(path_params.ClientIdParam())
	client.PUT("",
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, ids.UserIdFromContext,
		)),
	)
	client.PUT("/assertion-keys/:assertion_key_id",
		path_params.ClientAssertionKeyIdParam(),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientAssertionKeyHandler(
			p.commandBus, ids.ClientIdFromContext, ids.ClientAssertionKeyIdFromContext,
		)),
	)
	client.PUT("/secrets/:secret_id",
		path_params.ClientSecretIdParam(),
		ownership.Client(p.queryBus),
		p.executor.Run(securityHandlers.NewCreateClientSecretHandler(
			p.commandBus, ids.ClientIdFromContext, ids.ClientSecretIdFromContext,
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
	client.DELETE("/flows",
		ownership.Client(p.queryBus),
		p.executor.Run(clientHandlers.NewRemoveAllowedFlowHandler(p.commandBus, ids.ClientIdFromContext)),
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
