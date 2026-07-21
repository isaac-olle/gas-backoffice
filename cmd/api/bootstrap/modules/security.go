package modules

import (
	"github.com/golang-jwt/jwt/v5"
	clientAssertionKeyCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/create"
	clientAssertionKeyHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/hard_delete"
	clientVerifyAssertion "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/verify_assertion"
	clientAssertionKeyDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services"
	clientAssertionKeyInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/infrastructure"
	clientSecretCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/create"
	clientSecretHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/hard_delete"
	clientSecretVerify "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/verify_secret"
	clientSecretInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/infrastructure"
	keyPairEventHandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/_event_handlers"
	keyPairCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/create"
	generateJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt"
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	keyPairDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	userAssertionKeyCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys/application/create"
	verifyAssertion "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys/application/verify_assertion"
	userAssertionKeyDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys/domain_services"
	userAssertionKeyInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSecurity(pool *pgxpool.Pool, cmdBus command.Bus, qryBus query.Bus, evtBus event2.Bus, subsManager *SubscriptionManager, appName string) {
	keyPairRepo := keyPairInfra.NewPostgresKeyPairRepository(pool)
	clientAssertionKeyRepo := clientAssertionKeyInfra.NewPostgresClientAssertionKeyRepository(pool)
	clientAssertionKeyFinder := clientAssertionKeyDomainServices.NewClientAssertionKeyFinder(clientAssertionKeyRepo)
	userAssertionKeyRepo := userAssertionKeyInfra.NewPostgresUserAssertionKeyRepository(pool)
	clientSecretRepo := clientSecretInfra.NewPostgresClientSecretRepository(pool)

	jwtService := sharedJwt.NewJwtService(jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512, jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512)
	keyPairFinder := keyPairDomainServices.NewKeyPairFinder(keyPairRepo)

	createKeyPairService := keyPairCreate.NewService(keyPairRepo, evtBus)
	cmdBus.Register(keyPairCreate.CommandType, keyPairCreate.NewCommandHandler(createKeyPairService))
	cmdBus.Register(clientAssertionKeyCreate.CommandType, clientAssertionKeyCreate.NewCommandHandler(clientAssertionKeyCreate.NewService(clientAssertionKeyRepo, evtBus)))
	cmdBus.Register(clientAssertionKeyHardDelete.CommandType, clientAssertionKeyHardDelete.NewCommandHandler(clientAssertionKeyHardDelete.NewService(clientAssertionKeyRepo)))
	cmdBus.Register(userAssertionKeyCreate.CommandType, userAssertionKeyCreate.NewCommandHandler(userAssertionKeyCreate.NewService(userAssertionKeyRepo, evtBus)))
	clientSecretHasher := clientSecretInfra.NewArgon2SecretHasher()
	cmdBus.Register(clientSecretCreate.CommandType, clientSecretCreate.NewCommandHandler(clientSecretCreate.NewService(clientSecretRepo, clientSecretHasher, evtBus)))
	cmdBus.Register(clientSecretHardDelete.CommandType, clientSecretHardDelete.NewCommandHandler(clientSecretHardDelete.NewService(clientSecretRepo)))

	clientSecretVerifier := clientSecretInfra.NewArgon2SecretVerifier(clientSecretRepo)
	qryBus.Register(clientSecretVerify.QueryType, clientSecretVerify.NewQueryHandler(clientSecretVerify.NewService(clientSecretVerifier)))

	userAssertionKeyFinder := userAssertionKeyDomainServices.NewUserAssertionKeyFinder(userAssertionKeyRepo)
	userAssertionVerifier := userAssertionKeyInfra.NewJwtAssertionVerifier(jwtService)
	userJwtDecoder := userAssertionKeyInfra.NewJwtDecoder()
	qryBus.Register(verifyAssertion.QueryType, verifyAssertion.NewQueryHandler(verifyAssertion.NewService(userAssertionKeyFinder, userAssertionVerifier, userJwtDecoder, appName)))

	clientAssertionVerifier := clientAssertionKeyInfra.NewJwtAssertionVerifier(clientAssertionKeyFinder, jwtService)
	qryBus.Register(clientVerifyAssertion.QueryType, clientVerifyAssertion.NewQueryHandler(clientVerifyAssertion.NewService(clientAssertionVerifier)))
	keyPairTokenGenerator := keyPairInfra.NewJwtTokenGenerator(jwtService)
	keyPairTokenVerifier := keyPairInfra.NewJwtTokenVerifier(jwtService)
	keyPairJwtDecoder := keyPairInfra.NewJwtDecoder()
	qryBus.Register(generateJwt.QueryType, generateJwt.NewQueryHandler(generateJwt.NewService(keyPairFinder, keyPairTokenGenerator)))
	qryBus.Register(generateJwtByPrId.QueryType, generateJwtByPrId.NewQueryHandler(generateJwtByPrId.NewService(keyPairFinder, keyPairTokenGenerator)))
	qryBus.Register(verifyJwt.QueryType, verifyJwt.NewQueryHandler(verifyJwt.NewService(keyPairFinder, keyPairTokenVerifier, keyPairJwtDecoder, appName)))

	createKeyPairOnPRCreated := keyPairEventHandlers.NewCreateKeyPairOnProtectedResourceCreatedEventHandler(createKeyPairService)
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backofficeEvents.ProtectedResourceCreated](createKeyPairOnPRCreated.Handle))
}
