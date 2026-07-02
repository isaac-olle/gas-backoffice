package route_providers

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/server/handlers/admins"
	clientHandlers "github.com/isaac-olle/gas/internal/server/handlers/clients"
	"github.com/isaac-olle/gas/internal/server/handlers/helpers"
	prHandlers "github.com/isaac-olle/gas/internal/server/handlers/protected_resources"
	scopeHandlers "github.com/isaac-olle/gas/internal/server/handlers/scopes"
	tenantHandlers "github.com/isaac-olle/gas/internal/server/handlers/tenants"
	"github.com/isaac-olle/gas/internal/server/middlewares"
	"github.com/isaac-olle/gas/internal/server/middlewares/ownership"
	"github.com/isaac-olle/gas/internal/server/middlewares/path_params"
	"github.com/isaac-olle/gas/internal/server/middlewares/system_resource"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func registerAdmins(g *gin.RouterGroup, p *Backoffice) {
	a := g.Group("/admins")
	a.POST("/token",
		p.executor.Run(admins.NewTokenHandler(p.queryBus, p.signingKeyId, p.jwtIssuer)),
	)

	authenticated := a.Group("")
	authenticated.Use(
		middlewares.AccessToken(p.queryBus),
		middlewares.RoleMiddleware(vo.AdminRole),
	)
	authenticated.GET("",
		p.executor.Run(admins.NewGetAdminHandler(p.queryBus, helpers.AdminIdFromActor)),
	)

	tenant := authenticated.Group("/tenants/:tenant_id")
	tenant.Use(
		path_params.TenantIdParam(),
	)
	tenant.PUT("",
		p.executor.Run(tenantHandlers.NewCreateTenantHandler(
			p.commandBus, helpers.AdminIdFromActor, ids.TenantIdFromContext,
		)),
	)
	tenant.DELETE("",
		ownership.Tenant(p.queryBus),
		system_resource.Tenant(p.queryBus),
		p.executor.Run(tenantHandlers.NewSoftDeleteTenantHandler(p.commandBus, ids.TenantIdFromContext)),
	)
	tenant.GET("",
		ownership.Tenant(p.queryBus),
		p.executor.Run(tenantHandlers.NewGetTenantHandler(p.queryBus, ids.TenantIdFromContext)),
	)
	tenant.PATCH("/name",
		ownership.Tenant(p.queryBus),
		system_resource.Tenant(p.queryBus),
		p.executor.Run(tenantHandlers.NewChangeTenantNameHandler(p.commandBus, ids.TenantIdFromContext)),
	)

	pr := tenant.Group("/protected-resources/:protected_resource_id")
	pr.Use(
		path_params.ProtectedResourceIdParam(),
	)
	pr.PUT("",
		ownership.Tenant(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(prHandlers.NewCreateProtectedResourceHandler(
			p.commandBus,
			helpers.TenantOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.TenantIdFromContext,
		)),
		p.executor.Run(scopeHandlers.NewCreateScopesFromBodyHandler(
			p.commandBus,
			helpers.TenantOwnershipFromContext,
			ids.ProtectedResourceIdFromContext,
			ids.TenantIdFromContext,
			helpers.GenerateScopeId,
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
		p.executor.Run(prHandlers.NewUpdateProtectedResourceNameHandler(
			p.commandBus, ids.ProtectedResourceIdFromContext,
		)),
	)
	pr.DELETE("",
		ownership.ProtectedResource(p.queryBus),
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

	client := tenant.Group("/clients/:client_id")
	client.Use(
		path_params.ClientIdParam(),
	)
	client.PUT("",
		ownership.Tenant(p.queryBus),
		middlewares.Transactor(p.pgxPool),
		p.executor.Run(clientHandlers.NewCreateClientHandler(
			p.commandBus, ids.ClientIdFromContext, ids.TenantIdFromContext, helpers.TenantOwnershipFromContext,
		)),
		p.executor.Run(clientHandlers.NewCreateClientAssertionKeyHandler(
			p.commandBus, ids.ClientIdFromContext,
		)),
		p.executor.Run(clientHandlers.NewCreateClientSecretHandler(
			p.commandBus, ids.ClientIdFromContext,
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
