package modules

import (
	sessionEventHandlers "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/_event_handlers"
	buildAuthResponse "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/build_authorization_response"
	createAuthSession "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/create_authorization_session"
	exchangeCode "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/exchange_code"
	issueCC "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/issue_client_credentials"
	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/domain_services"
	sessionInfra "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/infrastructure"
	sessionPersistance "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/infrastructure/persistance"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	protocolEvents "github.com/isaac-olle/gas/internal/shared/events/protocol"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProtocol(pool *pgxpool.Pool, cmdBus command.Bus, qryBus query.Bus, evtBus event2.Bus, subsManager *SubscriptionManager, signingKeyId string, issuer string) {
	sessionRepo := sessionPersistance.NewPostgresAuthorizationCodeSessionRepository(pool)
	ccSessionRepo := sessionPersistance.NewPostgresClientCredentialsSessionRepository(pool)
	clientRepo := sessionInfra.NewPostgresClientRepository(pool)
	codeGenerator := sessionInfra.NewRandomCodeGenerator()
	clientFinder := domain_services.NewClientFinder(clientRepo)
	sessionFinder := domain_services.NewSessionFinder(sessionRepo)
	sessionTokenGenerator := sessionInfra.NewBusSessionTokenGenerator(qryBus, signingKeyId, issuer)

	cmdBus.Register(createAuthSession.CommandType, createAuthSession.NewCommandHandler(
		createAuthSession.NewService(sessionRepo, clientFinder, evtBus, codeGenerator),
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
