package modules

import (
	"time"

	sessionExchangeToken "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create_session"
	sessionInitiate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/initiate_login"
	sessionRefresh "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/refresh_session"
	authInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/helpers/oidc"
	userProjectionInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/user_projection"

	authLogout "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/logout"
	clientEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	clientAddAllowedGrantType "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_grant_type"
	clientAddRedirectUri "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
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
	eventLogList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs/application/list"
	eventLogInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs/infrastructure"
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
	subscriptionList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/list"
	subscriptionRemoveMember "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/remove_member"
	subscriptionUpdateName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/update_name"
	subscriptionInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/infrastructure"
	userChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
	userCreateTenant "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create_user"
	userGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/get"
	userHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/hard_delete"
	userList "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/list"
	userSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/soft_delete"
	userDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	userInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
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
	userRepo := userInfra.NewPostgresUserRepository(pool)
	subscriptionRepo := subscriptionInfra.NewPostgresRepository(pool)
	sessionRepo := authInfra.NewPostgresSessionRepository(pool)
	prRepo := prInfra.NewPostgresProtectedResourceRepository(pool)
	clientRepo := clientInfra.NewPostgresClientRepository(pool)
	scopeRepo := scopeInfra.NewPostgresScopeRepository(pool)
	clientProjectionRepo := scopeInfra.NewClientProjectionRepository(pool)
	idpRepo := idpInfra.NewPostgresIdentityProviderRepository(pool)
	errorLogRepo := errorLogInfra.NewPostgresErrorLogRepository(pool)
	eventLogRepo := eventLogInfra.NewPostgresEventLogRepository(pool)

	userFinder := userDomainServices.NewUserFinder(userRepo)
	prFinder := prDomainServices.NewProtectedResourceFinder(prRepo)
	scopeFinder := scopeDomainServices.NewScopeFinder(scopeRepo)
	clientProjectionFinder := scopeDomainServices.NewClientProjectionFinder(clientProjectionRepo)
	clientFinder := clientDomainServices.NewClientFinder(clientRepo)
	idpFinder := idpDomainServices.NewIdentityProviderFinder(idpRepo)

	clientDuplicateChecker := clientDomainServices.NewDuplicateNameChecker(clientRepo)
	prDuplicateChecker := prDomainServices.NewDuplicateNameChecker(prRepo)
	scopeDuplicateChecker := scopeDomainServices.NewDuplicateNameChecker(scopeRepo)

	// clients
	subscriptionAuthorizer := sharedDomainServices.NewSubscriptionAuthorizer()
	cmdBus.Register(clientCreate.CommandType, clientCreate.NewCommandHandler(clientCreate.NewService(clientDuplicateChecker, clientRepo, evtBus)))
	cmdBus.Register(clientAddAllowedGrantType.CommandType, clientAddAllowedGrantType.NewCommandHandler(clientAddAllowedGrantType.NewService(clientFinder, clientRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(clientRemoveAllowedGrantType.CommandType,
		clientRemoveAllowedGrantType.NewCommandHandler(clientRemoveAllowedGrantType.NewService(clientFinder,
			clientRepo,
			evtBus,
			subscriptionAuthorizer)))
	cmdBus.Register(clientAddRedirectUri.CommandType, clientAddRedirectUri.NewCommandHandler(clientAddRedirectUri.NewService(clientFinder, clientRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(clientSoftDelete.CommandType, clientSoftDelete.NewCommandHandler(clientSoftDelete.NewService(clientFinder, clientRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(clientHardDelete.CommandType, clientHardDelete.NewCommandHandler(clientHardDelete.NewService(clientFinder, clientRepo, evtBus, subscriptionAuthorizer)))

	// protected resources
	cmdBus.Register(prCreate.CommandType, prCreate.NewCommandHandler(prCreate.NewService(prDuplicateChecker, prRepo, evtBus)))
	cmdBus.Register(prUpdateName.CommandType, prUpdateName.NewCommandHandler(prUpdateName.NewService(prFinder, prDuplicateChecker, prRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(prSoftDelete.CommandType, prSoftDelete.NewCommandHandler(prSoftDelete.NewService(prFinder, prRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(prHardDelete.CommandType, prHardDelete.NewCommandHandler(prHardDelete.NewService(prFinder, prRepo, evtBus, subscriptionAuthorizer)))

	// scopes
	cmdBus.Register(scopeCreate.CommandType, scopeCreate.NewCommandHandler(scopeCreate.NewService(scopeDuplicateChecker, scopeRepo, evtBus)))
	cmdBus.Register(scopeAddAllowedClient.CommandType,
		scopeAddAllowedClient.NewCommandHandler(scopeAddAllowedClient.NewService(scopeFinder,
			scopeRepo,
			evtBus,
			clientProjectionFinder,
			subscriptionAuthorizer)))
	cmdBus.Register(scopeRemoveAllowedClient.CommandType,
		scopeRemoveAllowedClient.NewCommandHandler(scopeRemoveAllowedClient.NewService(scopeFinder,
			scopeRepo,
			evtBus,
			clientProjectionFinder,
			subscriptionAuthorizer)))
	cmdBus.Register(scopeSoftDelete.CommandType, scopeSoftDelete.NewCommandHandler(scopeSoftDelete.NewService(scopeFinder, scopeRepo, subscriptionAuthorizer)))
	cmdBus.Register(scopeHardDelete.CommandType, scopeHardDelete.NewCommandHandler(scopeHardDelete.NewService(scopeFinder, scopeRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(scopeEnable.CommandType, scopeEnable.NewCommandHandler(scopeEnable.NewService(scopeFinder, scopeRepo, evtBus, subscriptionAuthorizer)))
	cmdBus.Register(scopeDisable.CommandType, scopeDisable.NewCommandHandler(scopeDisable.NewService(scopeFinder, scopeRepo, evtBus, subscriptionAuthorizer)))

	// tenants
	cmdBus.Register(userCreateTenant.CommandType, userCreateTenant.NewCommandHandler(userCreateTenant.NewService(userRepo, evtBus)))
	cmdBus.Register(userChangeName.CommandType, userChangeName.NewCommandHandler(userChangeName.NewService(userFinder, userRepo, evtBus)))
	cmdBus.Register(userSoftDelete.CommandType, userSoftDelete.NewCommandHandler(userSoftDelete.NewService(userFinder, userRepo, evtBus)))
	cmdBus.Register(userHardDelete.CommandType, userHardDelete.NewCommandHandler(userHardDelete.NewService(userFinder, userRepo, evtBus)))

	// subscriptions
	cmdBus.Register(subscriptionCreate.CommandType, subscriptionCreate.NewCommandHandler(subscriptionCreate.NewService(subscriptionRepo, evtBus)))
	cmdBus.Register(subscriptionAddMember.CommandType, subscriptionAddMember.NewCommandHandler(subscriptionAddMember.NewService(subscriptionRepo, evtBus)))
	cmdBus.Register(subscriptionRemoveMember.CommandType, subscriptionRemoveMember.NewCommandHandler(subscriptionRemoveMember.NewService(subscriptionRepo, evtBus)))
	cmdBus.Register(subscriptionUpdateName.CommandType, subscriptionUpdateName.NewCommandHandler(subscriptionUpdateName.NewService(subscriptionRepo, evtBus)))
	cmdBus.Register(authLogout.CommandType, authLogout.NewCommandHandler(authLogout.NewService(sessionRepo)))

	// identity providers
	cmdBus.Register(idpCreate.CommandType, idpCreate.NewCommandHandler(idpCreate.NewService(idpRepo, evtBus, encrypter)))

	// error logs
	cmdBus.Register(errorLogCreate.CommandType, errorLogCreate.NewCommandHandler(errorLogCreate.NewService(errorLogRepo)))

	// get queries
	qryBus.Register(userGet.QueryType, userGet.NewQueryHandler(userGet.NewService(userFinder)))
	qryBus.Register(prGet.QueryType, prGet.NewQueryHandler(prGet.NewService(prFinder, subscriptionAuthorizer)))
	qryBus.Register(scopeGet.QueryType, scopeGet.NewQueryHandler(scopeGet.NewService(scopeFinder, subscriptionAuthorizer)))
	qryBus.Register(clientGet.QueryType, clientGet.NewQueryHandler(clientGet.NewService(clientFinder, subscriptionAuthorizer)))
	qryBus.Register(prCheckOwnership.QueryType, prCheckOwnership.NewQueryHandler(prCheckOwnership.NewService(prFinder)))
	qryBus.Register(idpGet.QueryType, idpGet.NewQueryHandler(idpGet.NewService(idpFinder)))
	qryBus.Register(userList.QueryType, userList.NewQueryHandler(userList.NewService(userRepo)))
	qryBus.Register(prList.QueryType, prList.NewQueryHandler(prList.NewService(prRepo, subscriptionAuthorizer)))
	qryBus.Register(errorLogList.QueryType, errorLogList.NewQueryHandler(errorLogList.NewService(errorLogRepo)))
	qryBus.Register(eventLogList.QueryType, eventLogList.NewQueryHandler(eventLogList.NewService(eventLogRepo)))
	qryBus.Register(clientList.QueryType, clientList.NewQueryHandler(clientList.NewService(clientRepo, subscriptionAuthorizer)))
	qryBus.Register(subscriptionList.QueryType, subscriptionList.NewQueryHandler(subscriptionList.NewService(subscriptionRepo)))

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
	qryBus.Register(sessionExchangeToken.QueryType,
		sessionExchangeToken.NewQueryHandler(sessionExchangeToken.NewService(tokenClient,
			stateEncoder,
			userProjectionRepo,
			tokenGenerator,
			tokenGenerator,
			sessionRepo,
			evtBus)))
	refreshTokenVerifier := authInfra.NewRefreshTokenVerifier(qryBus, backofficeProtectedResourceId)
	qryBus.Register(sessionRefresh.QueryType, sessionRefresh.NewQueryHandler(sessionRefresh.NewService(refreshTokenVerifier, tokenGenerator, tokenGenerator, sessionRepo, evtBus)))

}
