package modules

import (
	"time"

	userProjectionInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/_shared/user_projection/infrastructure"
	sessionExchangeToken "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create_session"
	sessionInitiate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/initiate_login"
	authInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/helpers/oidc"

	adminChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/change_name"
	adminCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/create"
	adminGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/get"
	adminDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	adminInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/infrastructure"
	clientEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	clientAddAllowedGrantType "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_grant_type"
	clientAddRedirectUri "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
	clientCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	clientGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/get"
	clientHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/hard_delete"
	clientList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/list"
	clientRemoveAllowedGrantType "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_grant_type"
	clientSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/soft_delete"
	clientDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/infrastructure"
	errorLogCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/create"
	errorLogList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/list"
	errorLogInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/infrastructure"
	idpCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/create"
	idpGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/get"
	idpDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/domain_services"
	idpInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/infrastructure"
	prEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/_event_handlers"
	prCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	prCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
	prGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/get"
	prHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/hard_delete"
	prList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/list"
	prSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/soft_delete"
	prUpdateName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/update_name"
	prDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/infrastructure"
	scopeEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/_event_handlers"
	scopeAddAllowedClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	scopeCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	scopeDisable "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/disable"
	scopeEnable "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/enable"
	scopeGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/get"
	scopeHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	scopeRemoveAllowedClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
	scopeSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/soft_delete"
	scopeDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopeInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/infrastructure"
	subscriptionEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/_event_handlers"
	subscriptionAddMember "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/add_member"
	subscriptionCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/create"
	subscriptionRemoveMember "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/remove_member"
	subscriptionInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/infrastructure"
	userChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/change_name"
	userCreateTenant "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/create_tenant"
	userGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/get"
	userHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/hard_delete"
	userList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/list"
	userSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/soft_delete"
	userDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	userInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBackoffice(
	pool *pgxpool.Pool, cmdBus command.Bus, qryBus query.Bus, evtBus event2.Bus,
	subsManager *SubscriptionManager, encrypter persistance.Encrypter,
	backofficeClientId, backofficeProtectedResourceId,
	backofficeClientRedirectUri string, backofficeClientSecret string, tokenEndpoint string,
	authorizeEndpoint string, jwksEndpoint string,
) {
	userRepo := userInfra.NewPostgresTenantRepository(pool)
	adminRepo := adminInfra.NewPostgresAdminRepository(pool)
	subscriptionRepo := subscriptionInfra.NewPostgresRepository(pool)
	prRepo := prInfra.NewPostgresProtectedResourceRepository(pool)
	clientRepo := clientInfra.NewPostgresClientRepository(pool)
	scopeRepo := scopeInfra.NewPostgresScopeRepository(pool)
	idpRepo := idpInfra.NewPostgresIdentityProviderRepository(pool)
	errorLogRepo := errorLogInfra.NewPostgresErrorLogRepository(pool)

	userFinder := userDomainServices.NewTenantFinder(userRepo)
	prFinder := prDomainServices.NewProtectedResourceFinder(prRepo)
	scopeFinder := scopeDomainServices.NewScopeFinder(scopeRepo)
	clientFinder := clientDomainServices.NewClientFinder(clientRepo)
	idpFinder := idpDomainServices.NewIdentityProviderFinder(idpRepo)
	adminFinder := adminDomainServices.NewAdminFinder(adminRepo)

	clientDuplicateChecker := clientDomainServices.NewDuplicateNameChecker(clientRepo)
	prDuplicateChecker := prDomainServices.NewDuplicateNameChecker(prRepo)
	scopeDuplicateChecker := scopeDomainServices.NewDuplicateNameChecker(scopeRepo)

	// clients
	cmdBus.Register(clientCreate.CommandType, clientCreate.NewCommandHandler(clientCreate.NewService(clientDuplicateChecker, clientRepo, evtBus)))
	cmdBus.Register(clientAddAllowedGrantType.CommandType, clientAddAllowedGrantType.NewCommandHandler(clientAddAllowedGrantType.NewService(clientFinder, clientRepo, evtBus)))
	cmdBus.Register(clientRemoveAllowedGrantType.CommandType, clientRemoveAllowedGrantType.NewCommandHandler(clientRemoveAllowedGrantType.NewService(clientFinder, clientRepo, evtBus)))
	cmdBus.Register(clientAddRedirectUri.CommandType, clientAddRedirectUri.NewCommandHandler(clientAddRedirectUri.NewService(clientFinder, clientRepo, evtBus)))
	cmdBus.Register(clientSoftDelete.CommandType, clientSoftDelete.NewCommandHandler(clientSoftDelete.NewService(clientFinder, clientRepo, evtBus)))
	cmdBus.Register(clientHardDelete.CommandType, clientHardDelete.NewCommandHandler(clientHardDelete.NewService(clientFinder, clientRepo)))

	// protected resources
	cmdBus.Register(prCreate.CommandType, prCreate.NewCommandHandler(prCreate.NewService(prDuplicateChecker, prRepo, evtBus)))
	cmdBus.Register(prUpdateName.CommandType, prUpdateName.NewCommandHandler(prUpdateName.NewService(prFinder, prDuplicateChecker, prRepo, evtBus)))
	cmdBus.Register(prSoftDelete.CommandType, prSoftDelete.NewCommandHandler(prSoftDelete.NewService(prFinder, prRepo, evtBus)))
	cmdBus.Register(prHardDelete.CommandType, prHardDelete.NewCommandHandler(prHardDelete.NewService(prFinder, prRepo, evtBus)))

	// scopes
	cmdBus.Register(scopeCreate.CommandType, scopeCreate.NewCommandHandler(scopeCreate.NewService(scopeDuplicateChecker, scopeRepo, evtBus)))
	cmdBus.Register(scopeAddAllowedClient.CommandType, scopeAddAllowedClient.NewCommandHandler(scopeAddAllowedClient.NewService(scopeFinder, scopeRepo, evtBus)))
	cmdBus.Register(scopeRemoveAllowedClient.CommandType, scopeRemoveAllowedClient.NewCommandHandler(scopeRemoveAllowedClient.NewService(scopeFinder, scopeRepo, evtBus)))
	cmdBus.Register(scopeSoftDelete.CommandType, scopeSoftDelete.NewCommandHandler(scopeSoftDelete.NewService(scopeFinder, scopeRepo)))
	cmdBus.Register(scopeHardDelete.CommandType, scopeHardDelete.NewCommandHandler(scopeHardDelete.NewService(scopeFinder, scopeRepo, evtBus)))
	cmdBus.Register(scopeEnable.CommandType, scopeEnable.NewCommandHandler(scopeEnable.NewService(scopeFinder, scopeRepo, evtBus)))
	cmdBus.Register(scopeDisable.CommandType, scopeDisable.NewCommandHandler(scopeDisable.NewService(scopeFinder, scopeRepo, evtBus)))

	// tenants
	cmdBus.Register(userCreateTenant.CommandType, userCreateTenant.NewCommandHandler(userCreateTenant.NewService(userRepo, evtBus)))
	cmdBus.Register(userChangeName.CommandType, userChangeName.NewCommandHandler(userChangeName.NewService(userFinder, userRepo, evtBus)))
	cmdBus.Register(userSoftDelete.CommandType, userSoftDelete.NewCommandHandler(userSoftDelete.NewService(userFinder, userRepo, evtBus)))
	cmdBus.Register(userHardDelete.CommandType, userHardDelete.NewCommandHandler(userHardDelete.NewService(userFinder, userRepo, evtBus)))

	// admins
	cmdBus.Register(adminCreate.CommandType, adminCreate.NewCommandHandler(adminCreate.NewService(adminRepo, evtBus)))
	cmdBus.Register(adminChangeName.CommandType, adminChangeName.NewCommandHandler(adminChangeName.NewService(adminFinder, adminRepo, evtBus)))

	// subscriptions
	cmdBus.Register(subscriptionCreate.CommandType, subscriptionCreate.NewCommandHandler(subscriptionCreate.NewService(subscriptionRepo, evtBus)))
	cmdBus.Register(subscriptionAddMember.CommandType, subscriptionAddMember.NewCommandHandler(subscriptionAddMember.NewService(subscriptionRepo, evtBus)))
	cmdBus.Register(subscriptionRemoveMember.CommandType, subscriptionRemoveMember.NewCommandHandler(subscriptionRemoveMember.NewService(subscriptionRepo, evtBus)))

	// identity providers
	cmdBus.Register(idpCreate.CommandType, idpCreate.NewCommandHandler(idpCreate.NewService(idpRepo, evtBus, encrypter)))

	// error logs
	cmdBus.Register(errorLogCreate.CommandType, errorLogCreate.NewCommandHandler(errorLogCreate.NewService(errorLogRepo)))

	// get queries
	qryBus.Register(userGet.QueryType, userGet.NewQueryHandler(userGet.NewService(userFinder)))
	qryBus.Register(prGet.QueryType, prGet.NewQueryHandler(prGet.NewService(prFinder)))
	qryBus.Register(scopeGet.QueryType, scopeGet.NewQueryHandler(scopeGet.NewService(scopeFinder)))
	qryBus.Register(clientGet.QueryType, clientGet.NewQueryHandler(clientGet.NewService(clientFinder)))
	qryBus.Register(clientCheckOwnership.QueryType, clientCheckOwnership.NewQueryHandler(clientCheckOwnership.NewService(clientFinder)))
	qryBus.Register(prCheckOwnership.QueryType, prCheckOwnership.NewQueryHandler(prCheckOwnership.NewService(prFinder)))
	qryBus.Register(idpGet.QueryType, idpGet.NewQueryHandler(idpGet.NewService(idpFinder)))
	qryBus.Register(userList.QueryType, userList.NewQueryHandler(userList.NewService(userRepo)))
	qryBus.Register(prList.QueryType, prList.NewQueryHandler(prList.NewService(prRepo)))
	qryBus.Register(errorLogList.QueryType, errorLogList.NewQueryHandler(errorLogList.NewService(errorLogRepo)))
	qryBus.Register(clientList.QueryType, clientList.NewQueryHandler(clientList.NewService(clientRepo)))
	qryBus.Register(adminGet.QueryType, adminGet.NewQueryHandler(adminGet.NewService(adminFinder)))

	// event handlers
	addKeyHandler := clientEventHandlers.NewAddKeyOnClientAssertionKeyCreatedEventHandler(clientFinder, clientRepo)
	addSecretHandler := clientEventHandlers.NewAddSecretOnClientSecretCreatedEventHandler(clientFinder, clientRepo)
	addAllowedScopeHandler := clientEventHandlers.NewAddAllowedScopeOnScopeClientAddedEventHandler(clientFinder, clientRepo)
	removeAllowedScopeHandler := clientEventHandlers.NewRemoveAllowedScopeOnScopeClientRemovedEventHandler(clientFinder, clientRepo)
	addScopeIdOnScopeCreated := prEventHandlers.NewAddScopeIdOnScopeCreatedEventHandler(prFinder, prRepo)
	removeScopeIdOnScopeDeleted := prEventHandlers.NewRemoveScopeIdOnScopeDeletedEventHandler(prFinder, prRepo)
	removeClientFromScopesOnClientDeleted := scopeEventHandlers.NewRemoveClientFromScopesOnClientDeletedEventHandler(scopeRepo, evtBus)
	disableScopesOnPRDeleted := scopeEventHandlers.NewDisableScopesOnProtectedResourceDeletedEventHandler(scopeRepo, evtBus)
	subscriptionMemberProjection := subscriptionEventHandlers.NewMemberProjectionHandler(pool)

	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[security.ClientAssertionKeyCreated](addKeyHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[security.ClientSecretCreated](addSecretHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeClientAdded](addAllowedScopeHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeClientRemoved](removeAllowedScopeHandler.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeCreated](addScopeIdOnScopeCreated.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ScopeDeleted](removeScopeIdOnScopeDeleted.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ClientDeleted](removeClientFromScopesOnClientDeleted.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.ProtectedResourceDeleted](disableScopesOnPRDeleted.Handle))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.SubscriptionMemberAdded](subscriptionMemberProjection.Add))
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backoffice.SubscriptionMemberRemoved](subscriptionMemberProjection.Remove))

	idpConfigRepo := authInfra.NewPostgresIdpConfigRepository(pool)

	codeVerifierGenerator := oidc.NewRandomCodeVerifierGenerator()
	stateEncoder := oidc.NewStateEncoder(10*time.Minute, encrypter)
	requestBuilder := oidc.NewRequestBuilder(codeVerifierGenerator, stateEncoder)
	authorizer := authInfra.NewOidcAuthorizer(
		requestBuilder, authorizeEndpoint, backofficeClientId, backofficeClientRedirectUri,
	)
	qryBus.Register(sessionInitiate.QueryType, sessionInitiate.NewQueryHandler(sessionInitiate.NewService(idpConfigRepo, authorizer)))

	idTokenVerifier := authInfra.NewGasIdTokenVerifier(jwksEndpoint, backofficeProtectedResourceId, backofficeClientId)
	tokenClient := authInfra.NewHttpTokenEndpointClient(tokenEndpoint, backofficeClientId, backofficeClientSecret, backofficeClientRedirectUri, idTokenVerifier)
	userProjectionRepo := userProjectionInfra.NewUserProjectionRepository(pool)
	tokenGenerator := authInfra.NewBusTokenGenerator(qryBus, backofficeProtectedResourceId)
	qryBus.Register(sessionExchangeToken.QueryType, sessionExchangeToken.NewQueryHandler(sessionExchangeToken.NewService(tokenClient, stateEncoder, userProjectionRepo, tokenGenerator, evtBus)))

	return
}
