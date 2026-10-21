package modules

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection/infrastructure"
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
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	keyPairDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSecurity(
	pool *pgxpool.Pool,
	cmdBus command.Bus,
	qryBus query.Bus,
	evtBus event2.Bus,
	subsManager *SubscriptionManager,
	encrypter persistance.Encrypter,
) {
	keyPairRepo := keyPairInfra.NewPostgresKeyPairRepository(pool, encrypter)
	clientAssertionKeyRepo := clientAssertionKeyInfra.NewPostgresClientAssertionKeyRepository(pool)
	clientAssertionKeyFinder := clientAssertionKeyDomainServices.NewClientAssertionKeyFinder(clientAssertionKeyRepo)
	clientSecretRepo := clientSecretInfra.NewPostgresClientSecretRepository(pool)

	clientFinder := infrastructure.NewClientProjectionRepository(pool)

	keyPairFinder := keyPairDomainServices.NewKeyPairFinder(keyPairRepo)

	keyPairGenerator := keyPairInfra.NewPEMKeyGenerator()
	createKeyPairService := keyPairCreate.NewService(keyPairRepo, evtBus, keyPairGenerator)
	cmdBus.Register(keyPairCreate.CommandType, keyPairCreate.NewCommandHandler(createKeyPairService))
	cmdBus.Register(clientAssertionKeyCreate.CommandType, clientAssertionKeyCreate.NewCommandHandler(clientAssertionKeyCreate.NewService(clientFinder, clientAssertionKeyRepo, evtBus)))
	cmdBus.Register(clientAssertionKeyHardDelete.CommandType, clientAssertionKeyHardDelete.NewCommandHandler(clientAssertionKeyHardDelete.NewService(clientAssertionKeyFinder, clientAssertionKeyRepo)))
	clientSecretHasher := clientSecretInfra.NewArgon2SecretHasher()
	cmdBus.Register(clientSecretCreate.CommandType, clientSecretCreate.NewCommandHandler(clientSecretCreate.NewService(clientFinder, clientSecretRepo, clientSecretHasher, evtBus)))
	cmdBus.Register(clientSecretHardDelete.CommandType, clientSecretHardDelete.NewCommandHandler(clientSecretHardDelete.NewService(clientSecretRepo)))

	clientSecretVerifier := clientSecretInfra.NewArgon2SecretVerifier(clientSecretRepo)
	qryBus.Register(clientSecretVerify.QueryType, clientSecretVerify.NewQueryHandler(clientSecretVerify.NewService(clientSecretVerifier)))

	clientAssertionVerifier := clientAssertionKeyInfra.NewJwtAssertionVerifier(clientAssertionKeyFinder)
	qryBus.Register(clientVerifyAssertion.QueryType, clientVerifyAssertion.NewQueryHandler(clientVerifyAssertion.NewService(clientAssertionVerifier)))
	keyPairTokenGenerator := keyPairInfra.NewJwtTokenGenerator()
	keyPairTokenVerifier := keyPairInfra.NewJwtTokenVerifier()
	keyPairJwtDecoder := keyPairInfra.NewJwtDecoder()
	qryBus.Register(generateJwtByPrId.QueryType, generateJwtByPrId.NewQueryHandler(generateJwtByPrId.NewService(keyPairFinder, keyPairTokenGenerator)))
	qryBus.Register(verifyJwt.QueryType, verifyJwt.NewQueryHandler(verifyJwt.NewService(keyPairFinder, keyPairTokenVerifier, keyPairJwtDecoder)))

	createKeyPairOnPRCreated := keyPairEventHandlers.NewCreateKeyPairOnProtectedResourceCreatedEventHandler(createKeyPairService)
	subsManager.AddInMemorySubscription(event2.NewInMemorySubscriber[backofficeEvents.ProtectedResourceCreated](createKeyPairOnPRCreated.Handle))
}
