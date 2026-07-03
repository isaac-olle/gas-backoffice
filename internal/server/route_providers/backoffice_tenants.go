package route_providers

import (
	"github.com/gin-gonic/gin"
	clientHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	tenantHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func registerTenants(g *gin.RouterGroup, p *Backoffice) {
	t := g.Group("/tenants")
	t.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.TenantRole),
	)

	t.GET("",
		ownership.Tenant(p.queryBus),
		p.executor.Run(tenantHandlers.NewGetTenantHandler(p.queryBus, helpers.TenantIdFromActor)),
	)
	t.PATCH("/name",
		ownership.Tenant(p.queryBus),
		system_resource.Tenant(p.queryBus),
		p.executor.Run(tenantHandlers.NewChangeTenantNameHandler(p.commandBus, helpers.TenantIdFromActor)),
	)
	t.DELETE("",
		ownership.Tenant(p.queryBus),
		system_resource.Tenant(p.queryBus),
		p.executor.Run(tenantHandlers.NewSoftDeleteTenantHandler(p.commandBus, helpers.TenantIdFromActor)),
	)

	pr := t.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.PUT("",
		ownership.Tenant(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(prHandlers.NewCreateProtectedResourceHandler(
			p.commandBus,
			helpers.TenantOwnershipFromActor,
			ids.ProtectedResourceIdFromContext,
			helpers.TenantIdFromActor,
		)),
		p.executor.Run(scopeHandlers.NewCreateScopesFromBodyHandler(
			p.commandBus,
			helpers.TenantOwnershipFromActor,
			ids.ProtectedResourceIdFromContext,
			helpers.TenantIdFromActor,
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
		ownership.Tenant(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, helpers.TenantIdFromActor, helpers.TenantOwnershipFromActor,
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
