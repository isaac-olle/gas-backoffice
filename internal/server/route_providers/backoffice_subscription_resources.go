package route_providers

import (
	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	securityHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/security"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// registerSubscriptionResources registers routes for resources owned by a
// subscription (clients, protected resources, scopes, and their security
// sub-resources). Every route is scoped by :subscription_id — the caller must
// be a member of it (checked by path_params.SubscriptionIdParam).
func registerSubscriptionResources(g *gin.RouterGroup, p *Backoffice) {
	s := g.Group("/subscriptions/:subscription_id")
	s.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins),
		middlewares.RoleMiddleware(vo.AdminRole, vo.TenantRole),
		path_params.SubscriptionIdParam(),
	)

	s.GET("/protected-resources",
		p.executor.Run(prHandlers.NewListProtectedResourcesHandler(p.queryBus, ids.SubscriptionIdFromContext)),
	)
	s.GET("/clients",
		p.executor.Run(clientHandlers.NewListClientsHandler(p.queryBus, ids.SubscriptionIdFromContext)),
	)

	pr := s.Group("/protected-resources/:protected_resource_id")
	pr.Use(path_params.ProtectedResourceIdParam())
	pr.PUT("",
		p.executor.Run(
			prHandlers.NewCreateProtectedResourceHandler(
				p.commandExecutor,
				ids.ProtectedResourceIdFromContext,
				ids.SubscriptionIdFromContext,
			)),
	)
	pr.GET("",
		p.executor.Run(prHandlers.NewGetProtectedResourceHandler(
			p.queryBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.PATCH("/name",
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandExecutor, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.DELETE("",
		p.executor.Run(prHandlers.NewSoftDeleteProtectedResourceHandler(
			p.commandExecutor, ids.ProtectedResourceIdFromContext,
		)),
	)

	scope := pr.Group("/scopes/:scope_id")
	scope.Use(path_params.ScopeIdParam())
	scope.PUT("",
		p.executor.Run(scopeHandlers.NewCreateScopeHandler(
			p.commandExecutor,
			p.queryBus,
			ids.ProtectedResourceIdFromContext,
			ids.ScopeIdFromContext,
		)),
	)
	scope.GET("",
		p.executor.Run(scopeHandlers.NewGetScopeHandler(
			p.queryBus, ids.ScopeIdFromContext, ids.ProtectedResourceIdFromContext,
		)),
	)
	scope.DELETE("",
		p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ProtectedResourceIdFromContext,
		)),
	)
	scope.PATCH("/enable",
		p.executor.Run(scopeHandlers.NewEnableScopeHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ProtectedResourceIdFromContext,
		)),
	)
	scope.PATCH("/disable",
		p.executor.Run(scopeHandlers.NewDisableScopeHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ProtectedResourceIdFromContext,
		)),
	)
	scope.POST("/allowed-clients/:client_id",
		path_params.ClientIdParam(),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(
			p.commandExecutor, p.queryBus, ids.ScopeIdFromContext, ids.ClientIdFromContext,
			ids.ProtectedResourceIdFromContext,
		)),
	)
	scope.DELETE("/allowed-clients/:client_id",
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ClientIdFromContext,
			ids.ProtectedResourceIdFromContext,
		)),
	)

	client := s.Group("/clients/:client_id")
	client.Use(path_params.ClientIdParam())
	client.PUT("",
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.SubscriptionIdFromContext,
		)),
	)
	client.PUT("/assertion-keys/:assertion_key_id",
		path_params.ClientAssertionKeyIdParam(),
		p.executor.Run(securityHandlers.NewCreateClientAssertionKeyHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.ClientAssertionKeyIdFromContext,
		)),
	)
	client.PUT("/secrets/:secret_id",
		path_params.ClientSecretIdParam(),
		p.executor.Run(securityHandlers.NewCreateClientSecretHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.ClientSecretIdFromContext,
		)),
	)
	client.DELETE("",
		p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(
			p.commandExecutor, ids.ClientIdFromContext,
		)),
	)
	client.POST("/grant-types",
		p.executor.Run(clientHandlers.NewAddAllowedGrantTypeHandler(
			p.commandExecutor, ids.ClientIdFromContext,
		)),
	)
	client.DELETE("/grant-types",
		p.executor.Run(clientHandlers.NewRemoveAllowedGrantTypeHandler(
			p.commandExecutor, ids.ClientIdFromContext,
		)),
	)
	client.GET("",
		p.executor.Run(clientHandlers.NewGetClientHandler(
			p.queryBus, ids.ClientIdFromContext,
		)),
	)
	client.POST("/redirect-uris",
		p.executor.Run(clientHandlers.NewAddRedirectUriHandler(
			p.commandExecutor, ids.ClientIdFromContext,
		)),
	)
}
