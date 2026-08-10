package modules

import (
	sessionCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create"
	sessionGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/get"
	authDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/domain_services"
	authInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	clientEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	clientAddAllowedFlow "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_flow"
	clientAddRedirectUri "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
	clientCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	clientGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/get"
	clientHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/hard_delete"
	clientRemoveAllowedFlow "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_flow"
	clientSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/soft_delete"
	clientDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/infrastructure"
	idpGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/get"
	idpDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/domain_services"
	idpInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/infrastructure"
	prEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/_event_handlers"
	prCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	prCheckSystem "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_system_resource"
	prCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
	prGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/get"
	prHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/hard_delete"
	prSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/soft_delete"
	prUpdateName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/update_name"
	prDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/infrastructure"
	scopeAddAllowedClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	scopeCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/check_ownership"
	scopeCheckSystem "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/check_system_resource"
	scopeCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	scopeGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/get"
	scopeHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	scopeRemoveAllowedClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
	scopeSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/soft_delete"
	scopeDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopeInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/infrastructure"
	userEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/_event_handlers"
	userChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
	userCheckSystem "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_system_resource"
	userCheckTenantAccess "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_tenant_access"
	userCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create"
	userGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/get"
	userHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/hard_delete"
	userSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/soft_delete"
	userDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	userInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
	persistanceInfra "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBackoffice(pool *pgxpool.Pool, cmdBus command.Bus, qryBus query.Bus, evtBus event2.Bus, subsManager *SubscriptionManager, openIDScopeId string) {
	userRepo := persistanceInfra.NewInMemoryContextDecorator(userInfra.NewPostgresUserRepository(pool))
	prRepo := persistanceInfra.NewInMemoryContextDecorator(prInfra.NewPostgresProtectedResourceRepository(pool))
	clientRepo := persistanceInfra.NewInMemoryContextDecorator(clientInfra.NewPostgresClientRepository(pool))
	scopeRepo := persistanceInfra.NewInMemoryContextDecorator(scopeInfra.NewPostgresScopeRepository(pool))
	sessionRepo := persistanceInfra.NewInMemoryContextDecorator(authInfra.NewPostgresSessionRepository(pool))
	idpRepo := idpInfra.NewPostgresIdentityProviderRepository(pool)

	userFinder := userDomainServices.NewUserFinder(userRepo)
	prFinder := prDomainServices.NewProtectedResourceFinder(prRepo)
	scopeFinder := scopeDomainServices.NewScopeFinder(scopeRepo)
	clientFinder := clientDomainServices.NewClientFinder(clientRepo)
	sessionFinder := authDomainServices.NewSessionFinder(sessionRepo)
	idpFinder := idpDomainServices.NewIdentityProviderFinder(idpRepo)

	// clients
	cmdBus.Register(clientCreate.CommandType, clientCreate.NewCommandHandler(clientCreate.NewService(clientRepo, evtBus)))
	cmdBus.Register(clientAddAllowedFlow.CommandType, clientAddAllowedFlow.NewCommandHandler(clientAddAllowedFlow.NewService(clientRepo, evtBus)))
	cmdBus.Register(clientRemoveAllowedFlow.CommandType, clientRemoveAllowedFlow.NewCommandHandler(clientRemoveAllowedFlow.NewService(clientFinder, clientRepo, evtBus)))
	cmdBus.Register(clientAddRedirectUri.CommandType, clientAddRedirectUri.NewCommandHandler(clientAddRedirectUri.NewService(clientRepo, evtBus)))
	cmdBus.Register(clientSoftDelete.CommandType, clientSoftDelete.NewCommandHandler(clientSoftDelete.NewService(clientFinder, clientRepo, evtBus)))
	cmdBus.Register(clientHardDelete.CommandType, clientHardDelete.NewCommandHandler(clientHardDelete.NewService(clientRepo)))

	// protected resources
	cmdBus.Register(prCreate.CommandType, prCreate.NewCommandHandler(prCreate.NewService(prRepo, evtBus)))
	cmdBus.Register(prUpdateName.CommandType, prUpdateName.NewCommandHandler(prUpdateName.NewService(prFinder, prRepo, evtBus)))
	cmdBus.Register(prSoftDelete.CommandType, prSoftDelete.NewCommandHandler(prSoftDelete.NewService(prFinder, prRepo, evtBus)))
	cmdBus.Register(prHardDelete.CommandType, prHardDelete.NewCommandHandler(prHardDelete.NewService(prFinder, prRepo, evtBus)))

	// scopes
	cmdBus.Register(scopeCreate.CommandType, scopeCreate.NewCommandHandler(scopeCreate.NewService(scopeRepo, evtBus)))
	cmdBus.Register(scopeAddAllowedClient.CommandType, scopeAddAllowedClient.NewCommandHandler(scopeAddAllowedClient.NewService(scopeRepo, evtBus)))
	cmdBus.Register(scopeRemoveAllowedClient.CommandType, scopeRemoveAllowedClient.NewCommandHandler(scopeRemoveAllowedClient.NewService(scopeRepo, evtBus)))
	cmdBus.Register(scopeSoftDelete.CommandType, scopeSoftDelete.NewCommandHandler(scopeSoftDelete.NewService(scopeRepo)))
	cmdBus.Register(scopeHardDelete.CommandType, scopeHardDelete.NewCommandHandler(scopeHardDelete.NewService(scopeFinder, scopeRepo, evtBus)))

	// users
	cmdBus.Register(userCreate.CommandType, userCreate.NewCommandHandler(userCreate.NewService(userRepo, evtBus)))
	cmdBus.Register(userChangeName.CommandType, userChangeName.NewCommandHandler(userChangeName.NewService(userFinder, userRepo, evtBus)))
	cmdBus.Register(userSoftDelete.CommandType, userSoftDelete.NewCommandHandler(userSoftDelete.NewService(userFinder, userRepo, evtBus)))
	cmdBus.Register(userHardDelete.CommandType, userHardDelete.NewCommandHandler(userHardDelete.NewService(userFinder, userRepo, evtBus)))

	// sessions
	cmdBus.Register(sessionCreate.CommandType, sessionCreate.NewCommandHandler(sessionCreate.NewService(sessionRepo, evtBus)))

	// check_ownership queries
	qryBus.Register(userCheckTenantAccess.QueryType, userCheckTenantAccess.NewQueryHandler(userCheckTenantAccess.NewService(userFinder)))
	qryBus.Register(prCheckOwnership.QueryType, prCheckOwnership.NewQueryHandler(prCheckOwnership.NewService(prRepo)))
	qryBus.Register(scopeCheckOwnership.QueryType, scopeCheckOwnership.NewQueryHandler(scopeCheckOwnership.NewService(scopeRepo)))
	qryBus.Register(clientCheckOwnership.QueryType, clientCheckOwnership.NewQueryHandler(clientCheckOwnership.NewService(clientRepo)))

	// check_system_resource queries
	qryBus.Register(userCheckSystem.QueryType, userCheckSystem.NewQueryHandler(userCheckSystem.NewService(userFinder)))
	qryBus.Register(prCheckSystem.QueryType, prCheckSystem.NewQueryHandler(prCheckSystem.NewService(prFinder)))
	qryBus.Register(scopeCheckSystem.QueryType, scopeCheckSystem.NewQueryHandler(scopeCheckSystem.NewService(scopeFinder)))

	// get queries
	qryBus.Register(userGet.QueryType, userGet.NewQueryHandler(userGet.NewService(userFinder)))
	qryBus.Register(prGet.QueryType, prGet.NewQueryHandler(prGet.NewService(prFinder)))
	qryBus.Register(scopeGet.QueryType, scopeGet.NewQueryHandler(scopeGet.NewService(scopeFinder)))
	qryBus.Register(clientGet.QueryType, clientGet.NewQueryHandler(clientGet.NewService(clientFinder)))
	qryBus.Register(sessionGet.QueryType, sessionGet.NewQueryHandler(sessionGet.NewService(sessionFinder)))
	qryBus.Register(idpGet.QueryType, idpGet.NewQueryHandler(idpGet.NewService(idpFinder)))

	// event handlers
	addClientHandler := userEventHandlers.NewAddClientOnClientCreatedEventHandler(userRepo, evtBus)
	removeClientHandler := userEventHandlers.NewRemoveClientOnClientDeletedEventHandler(userRepo, evtBus)
	addPRHandler := userEventHandlers.NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(userRepo, evtBus)
	removePRHandler := userEventHandlers.NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(userRepo, evtBus)
	addScopeHandler := prEventHandlers.NewAddScopeOnScopeCreatedEventHandler(prRepo, evtBus)
	removeScopeHandler := prEventHandlers.NewRemoveScopeOnScopeDeletedEventHandler(prRepo, evtBus)
	addScopeToClientHandler := clientEventHandlers.NewAddScopeOnScopeClientAddedEventHandler(clientRepo, evtBus)
	removeScopeFromClientHandler := clientEventHandlers.NewRemoveScopeOnScopeClientRemovedEventHandler(clientRepo, evtBus)
	addKeyHandler := clientEventHandlers.NewAddKeyOnClientAssertionKeyCreatedEventHandler(clientFinder, clientRepo, evtBus)
	addSecretHandler := clientEventHandlers.NewAddSecretOnClientSecretCreatedEventHandler(clientFinder, clientRepo, evtBus)
	addOpenIDScopeHandler := clientEventHandlers.NewAddOpenIDScopeOnClientCreatedEventHandler(clientRepo, evtBus, func() (ids2.ScopeId, error) {
		return ids2.NewScopeIdFromString(openIDScopeId)
	})

	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ClientCreated](addClientHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ClientCreated](addOpenIDScopeHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ClientDeleted](removeClientHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ProtectedResourceCreated](addPRHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ProtectedResourceDeleted](removePRHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeCreated](addScopeHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeDeleted](removeScopeHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeClientAdded](addScopeToClientHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeClientRemoved](removeScopeFromClientHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[security.ClientAssertionKeyCreated](addKeyHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[security.ClientSecretCreated](addSecretHandler.Handle))
}
