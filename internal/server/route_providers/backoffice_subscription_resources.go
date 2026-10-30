package route_providers

import (
	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	securityHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/security"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/subscription"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// registerSubscriptionResources keeps subscription_id only where it is part
// of the creation/listing context. Existing resources are addressed by their
// own aggregate path and are authorized by OwnershipMiddleware.
func registerSubscriptionResources(g *gin.RouterGroup, p *Backoffice) {
	resources := g.Group("")
	resources.Use(
		middlewares.Authentication(p.queryBus, p.jwtIssuer, p.cookieOrigins, p.sessionRepository),
		middlewares.RoleMiddleware(vo.SystemAdministratorRole, vo.SubscriptionUserRole),
	)
	resources.GET("/protected-resources", p.executor.Run(
		prHandlers.NewListProtectedResourcesHandler(p.queryBus),
	))
	resources.GET("/clients", p.executor.Run(
		clientHandlers.NewListClientsHandler(p.queryBus),
	))

	// Collection and creation routes carry the target subscription explicitly.
	s := resources.Group("/subscriptions/:subscription_id")
	s.Use(middlewares.SubscriptionMiddleware(subscription.NewPathResolver("subscription_id")))
	prCreate := s.Group("/protected-resources/:protected_resource_id")
	prCreate.Use(path_params.ProtectedResourceIdParam())
	prCreate.PUT("", p.executor.Run(prHandlers.NewCreateProtectedResourceHandler(
		p.commandExecutor, ids.ProtectedResourceIdFromContext, ids.SubscriptionIdFromContext,
	)))

	scopeCreate := s.Group("/protected-resources/:protected_resource_id/scopes/:scope_id")
	scopeCreate.Use(path_params.ProtectedResourceIdParam(), path_params.ScopeIdParam())
	scopeCreate.PUT("", p.executor.Run(scopeHandlers.NewCreateScopeHandler(
		p.commandExecutor, p.queryBus, ids.ProtectedResourceIdFromContext, ids.ScopeIdFromContext,
	)))

	clientCreate := s.Group("/clients/:client_id")
	clientCreate.Use(path_params.ClientIdParam())
	clientCreate.PUT("", p.executor.Run(clientHandlers.NewCreateClientHandler(
		p.commandExecutor, ids.ClientIdFromContext, ids.SubscriptionIdFromContext,
	)))

	clientChildren := s.Group("/clients/:client_id")
	clientChildren.Use(
		path_params.ClientIdParam(),
	)
	clientChildren.PUT("/assertion-keys/:assertion_key_id", path_params.ClientAssertionKeyIdParam(),
		p.executor.Run(securityHandlers.NewCreateClientAssertionKeyHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.ClientAssertionKeyIdFromContext,
		)))
	clientChildren.PUT("/secrets/:secret_id", path_params.ClientSecretIdParam(),
		p.executor.Run(securityHandlers.NewCreateClientSecretHandler(
			p.commandExecutor, ids.ClientIdFromContext, ids.ClientSecretIdFromContext,
		)))

	// Existing resources are addressed without subscription_id. Ownership
	// resolves their real subscription and checks it against the token.
	pr := resources.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.GET("", p.executor.Run(prHandlers.NewGetProtectedResourceHandler(
		p.queryBus, ids.ProtectedResourceIdFromContext, helpers.AuthorFromContext,
	)))
	pr.PATCH("/name", p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
		p.commandExecutor, ids.ProtectedResourceIdFromContext,
	)))
	pr.DELETE("", p.executor.Run(prHandlers.NewSoftDeleteProtectedResourceHandler(
		p.commandExecutor, ids.ProtectedResourceIdFromContext,
	)))

	scope := resources.Group("/scopes/:scope_id")
	scope.Use(
		path_params.ScopeIdParam(),
	)
	scope.GET("", p.executor.Run(scopeHandlers.NewGetScopeHandler(
		p.queryBus, ids.ScopeIdFromContext, helpers.AuthorFromContext,
	)))
	scope.DELETE("", p.executor.Run(scopeHandlers.NewSoftDeleteScopeHandler(
		p.commandExecutor, ids.ScopeIdFromContext,
	)))
	scope.PATCH("/enable", p.executor.Run(scopeHandlers.NewEnableScopeHandler(
		p.commandExecutor, ids.ScopeIdFromContext,
	)))
	scope.PATCH("/disable", p.executor.Run(scopeHandlers.NewDisableScopeHandler(
		p.commandExecutor, ids.ScopeIdFromContext,
	)))
	scope.POST("/allowed-clients/:client_id", path_params.ClientIdParam(),
		p.executor.Run(scopeHandlers.NewAddScopeAllowedClientHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)))
	scope.DELETE("/allowed-clients/:client_id", path_params.ClientIdParam(),
		p.executor.Run(scopeHandlers.NewRemoveScopeAllowedClientHandler(
			p.commandExecutor, ids.ScopeIdFromContext, ids.ClientIdFromContext,
		)))

	client := resources.Group("/clients/:client_id")
	client.Use(path_params.ClientIdParam())
	client.GET("", p.executor.Run(clientHandlers.NewGetClientHandler(
		p.queryBus, ids.ClientIdFromContext, helpers.AuthorFromContext,
	)))
	client.DELETE("", p.executor.Run(clientHandlers.NewSoftDeleteClientHandler(
		p.commandExecutor, ids.ClientIdFromContext,
	)))
	client.POST("/grant-types", p.executor.Run(clientHandlers.NewAddAllowedGrantTypeHandler(
		p.commandExecutor, ids.ClientIdFromContext,
	)))
	client.DELETE("/grant-types", p.executor.Run(clientHandlers.NewRemoveAllowedGrantTypeHandler(
		p.commandExecutor, ids.ClientIdFromContext,
	)))
	client.POST("/redirect-uris", p.executor.Run(clientHandlers.NewAddRedirectUriHandler(
		p.commandExecutor, ids.ClientIdFromContext,
	)))
}
