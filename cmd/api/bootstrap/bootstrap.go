package bootstrap

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_managers"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_presenters"
	commandInfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/command/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	outboxapp "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/application/outbox"
	queryInfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/logger"
	persistanceInfra "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/workers"
	"strconv"
	"time"

	"github.com/isaac-olle/gas-backoffice/cmd/api/bootstrap/modules"
	"github.com/isaac-olle/gas-backoffice/config"
	"github.com/isaac-olle/gas-backoffice/internal/server"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/route_providers"
	auditapp "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/application/audit"
	eventinfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log, err := logger.NewZapLogger(cfg.Environment)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	postgresPool, err := pgxpool.New(context.Background(), cfg.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer postgresPool.Close()

	workerPool := workers.NewPool(10)
	outboxRepo := eventinfra.NewOutboxRepository(postgresPool)
	outboxEventBus := outboxapp.NewEventBus(outboxRepo, log)
	subsManager := modules.NewSubscriptionManager()
	uow := persistanceInfra.NewPostgresUnitOfWork(postgresPool)
	cmdBus := commandInfra.NewTransactionalCommandBus(commandInfra.NewCommandBus(), uow)
	commandExecutor := commandInfra.NewCommandExecutor(cmdBus, uow)
	qryBus := queryInfra.NewQueryBus()

	auditRepo := eventinfra.NewPostgresAuditRepository(postgresPool)
	subsManager.AddInMemorySubscription(auditapp.NewSubscriber(auditRepo))

	encrypter, err := persistanceInfra.NewAESEncrypter(cfg.Security.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to create encrypter: %w", err)
	}

	authorizeService, tokenService, authClientService, createSessionService, jwksService, _ := modules.RegisterProtocol(
		postgresPool,
		qryBus,
		outboxEventBus,
		subsManager,
		encrypter,
		uow,
	)
	tokenEndpoint := cfg.Application.ProtocolUrl + "/oidc/token"
	authorizeEndpoint := cfg.Application.ProtocolUrl + "/oidc/authorize"
	jwksEndpoint := cfg.Application.ProtocolUrl + "/oidc/jwks"
	modules.RegisterBackoffice(
		postgresPool,
		cmdBus,
		qryBus,
		outboxEventBus,
		subsManager,
		encrypter,
		cfg.Application.BackofficeClientId,
		cfg.Application.BackofficeProtectedResourceId,
		cfg.Application.BackofficeClientRedirectUri,
		cfg.Application.BackofficeClientSecret,
		tokenEndpoint,
		authorizeEndpoint,
		jwksEndpoint,
	)
	modules.RegisterSecurity(
		postgresPool,
		cmdBus,
		qryBus,
		outboxEventBus,
		subsManager,
		encrypter,
	)

	eventBus := subsManager.Build(log, workerPool)

	registry := event.NewBuilderRegistry()
	modules.RegisterEventBuilders(registry)

	port, err := strconv.Atoi(cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("invalid server port: %w", err)
	}

	presenter := error_presenters.NewHttpErrorPresenter()
	managers := []app_errors.ErrorManager{
		error_managers.NewLogErrorPresenter(log),
		error_managers.NewPostgresErrorManager(commandExecutor),
	}
	hm := middlewares.NewHandlerMiddleware()

	ctx, srv := server.New(
		context.Background(),
		port,
		cfg.Server.ShutdownTimeout(),
		cmdBus,
		qryBus,
		log,
		presenter,
		managers,
		cfg.Server.AllowedOrigins,
		route_providers.NewBackoffice(hm, qryBus, commandExecutor, postgresPool, cfg.Application.BackofficeProtectedResourceId, cfg.Server.CookieOrigins),
		route_providers.NewProtocol(hm, authorizeService, tokenService, authClientService, jwksService, createSessionService, managers...),
		route_providers.NewSeed(postgresPool, commandExecutor, cfg.Application.SeedToken),
	)

	consumer := eventinfra.NewOutboxConsumer(outboxRepo, eventBus, registry, 500*time.Millisecond, log)
	go consumer.Run(ctx)

	return srv.Run(ctx, log)
}
