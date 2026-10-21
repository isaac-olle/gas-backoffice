package modules

import (
	"time"

	protocolSharedDS "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/domain_services"
	protocolSharedInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/infrastructure"
	getJwks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/application/get_jwks"
	jwksDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/domain_services"
	jwksInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authorize"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/create_session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token/flows"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	sessionInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	sessionHelpers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/helpers"
	sessionPersistance "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/persistance"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/jackc/pgx/v5/pgxpool"
)

// idpStateTTL bounds how long a client has to complete the Google login
// flow before the encrypted idp_state sent to /oidc/authorize expires.
const idpStateTTL = 10 * time.Minute

func RegisterProtocol(
	pool *pgxpool.Pool, qryBus query.Bus, evtBus event2.Bus,
	subsManager *SubscriptionManager, encrypter persistance.Encrypter, uow persistance.UnitOfWork,
) (*authorize.Service, *token.Service, *authenticate_client.Service, *create_session.Service, *getJwks.Service, *domain_services.AuthorizationCodeSessionFinder) {
	sessionRepo := sessionPersistance.NewRequestCachedSessionRepository(
		sessionPersistance.NewPostgresAuthorizationCodeSessionRepository(pool, encrypter),
	)
	ccSessionRepo := sessionPersistance.NewPostgresClientCredentialsSessionRepository(pool)
	clientRepo := sessionInfra.NewPostgresClientRepository(pool)
	idpRepo := protocolSharedInfra.NewPostgresIdentityProviderRepository(pool, encrypter)
	idpFinder := protocolSharedDS.NewIdentityProviderFinder(idpRepo)
	sessionIdpClient := sessionInfra.NewAuthorizationURLBuilder()
	oidcTokenVerifier := sessionHelpers.NewOidcTokenVerifier()
	idpClient := sessionInfra.NewOidcIdentityProviderClient(oidcTokenVerifier)
	idpStateCodec := sessionInfra.NewIdpStateCodec(idpStateTTL, encrypter)
	codeGenerator := sessionInfra.NewRandomCodeGenerator()
	clientFinder := domain_services.NewClientFinder(clientRepo)
	authSessionFinder := domain_services.NewAuthorizationCodeSessionFinder(sessionRepo)
	ccSessionFinder := domain_services.NewClientCredentialsSessionFinder(ccSessionRepo)
	accessTokenGenerator := sessionInfra.NewBusAccessTokenGenerator(qryBus)
	idTokenGenerator := sessionInfra.NewBusIDTokenGenerator(qryBus)
	authCodeTokenResponseGenerator := domain_services.NewAuthorizationCodeTokenResponseGenerator(accessTokenGenerator, idTokenGenerator)
	ccTokenResponseGenerator := domain_services.NewClientCredentialsTokenResponseGenerator(accessTokenGenerator)

	authorizeService := authorize.NewService(
		clientFinder,
		idpFinder,
		sessionIdpClient,
		idpStateCodec,
	)

	tokenService := token.NewService(
		flows.NewAuthorizationCodeGrantHandler(
			authSessionFinder,
			clientFinder,
			sessionRepo,
			evtBus,
			authCodeTokenResponseGenerator,
			uow,
		),
		flows.NewClientCredentialsGrantHandler(
			clientFinder,
			ccSessionRepo,
			evtBus,
			ccSessionFinder,
			ccTokenResponseGenerator,
			uow,
		),
	)

	authClientService := authenticate_client.NewService(
		authenticate_client.NewBasicMethod(sessionInfra.NewBusClientSecretVerifier(qryBus)),
		authenticate_client.NewJWTAssertionMethod(sessionInfra.NewBusClientAssertionVerifier(qryBus)),
		authenticate_client.NewBodyCredentialsMethod(sessionInfra.NewBusClientSecretVerifier(qryBus)),
		authenticate_client.NewPublicClientMethod(clientFinder),
	)

	createSessionService := create_session.NewService(
		idpStateCodec,
		clientFinder,
		idpFinder,
		idpClient,
		sessionRepo,
		codeGenerator,
		evtBus,
	)

	jwksService := getJwks.NewService(jwksDomainServices.NewSigningKeyFinder(jwksInfra.NewPostgresSigningKeyRepository(pool)))

	return authorizeService, tokenService, authClientService, createSessionService, jwksService, authSessionFinder
}
