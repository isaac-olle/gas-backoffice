package modules

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	sessionEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/_event_handlers"
	authenticateWithIdp "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_with_idp"
	buildAuthResponse "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/build_authorization_response"
	createAuthSession "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_authorization_session"
	createLoginSession "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_login_authorization_session"
	exchangeCode "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/exchange_code"
	issueCC "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/issue_client_credentials"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	sessionInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	sessionPersistance "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	protocolEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BackofficeLoginConfig identifies the backoffice frontend as a gas OAuth
// client for the federated-login session (create_login_authorization_session).
type BackofficeLoginConfig struct {
	ClientId            string
	RedirectUri         string
	Scope               string
	ProtectedResourceId string
}

func RegisterProtocol(pool *pgxpool.Pool, cmdBus command.Bus, qryBus query.Bus, evtBus event2.Bus, subsManager *SubscriptionManager, issuer string, login BackofficeLoginConfig, publicBaseURL string) {
	sessionRepo := sessionPersistance.NewPostgresAuthorizationCodeSessionRepository(pool)
	ccSessionRepo := sessionPersistance.NewPostgresClientCredentialsSessionRepository(pool)
	clientRepo := sessionInfra.NewPostgresClientRepository(pool)
	codeGenerator := sessionInfra.NewRandomCodeGenerator()
	clientFinder := domain_services.NewClientFinder(clientRepo)
	sessionFinder := domain_services.NewSessionFinder(sessionRepo)
	sessionTokenGenerator := sessionInfra.NewBusSessionTokenGenerator(qryBus, issuer)

	cmdBus.Register(createAuthSession.CommandType, createAuthSession.NewCommandHandler(
		createAuthSession.NewService(sessionRepo, clientFinder, evtBus, codeGenerator),
	))

	cmdBus.Register(createLoginSession.CommandType, createLoginSession.NewCommandHandler(
		createLoginSession.NewService(
			sessionRepo, evtBus, codeGenerator,
			login.ClientId, login.RedirectUri, login.Scope, login.ProtectedResourceId,
		),
	))

	// Federated login (gas as OIDC relying party). Only Google is registered.
	idpFinder := domain_services.NewIdentityProviderFinder(sessionInfra.NewPostgresIdentityProviderRepository(pool))
	idpClientFactory := domain_services.NewIdentityProviderClientFactory(map[string]session.IdentityProviderClient{
		"google": sessionInfra.NewGoogleIdentityProviderClient(publicBaseURL),
	})
	qryBus.Register(authenticateWithIdp.QueryType, authenticateWithIdp.NewQueryHandler(
		authenticateWithIdp.NewService(sessionFinder, idpFinder, idpClientFactory),
	))

	qryBus.Register(buildAuthResponse.QueryType, buildAuthResponse.NewQueryHandler(
		buildAuthResponse.NewService(sessionFinder),
	))

	qryBus.Register(exchangeCode.QueryType, exchangeCode.NewQueryHandler(
		exchangeCode.NewService(sessionFinder, clientFinder, sessionRepo, sessionTokenGenerator),
	))

	qryBus.Register(issueCC.QueryType, issueCC.NewQueryHandler(
		issueCC.NewService(clientFinder, sessionTokenGenerator, evtBus),
	))

	persistCCSession := sessionEventHandlers.NewPersistClientCredentialsSessionOnCreatedEventHandler(ccSessionRepo)
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[protocolEvents.ClientCredentialsSessionCreated](persistCCSession.Handle))
}
