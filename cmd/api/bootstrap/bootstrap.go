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

	if err = applyMigrations(log, cfg.Databases.MigrationsPath, cfg.Databases.MigrationsDSN); err != nil {
		return err
	}

	postgresPool, err := pgxpool.New(context.Background(), cfg.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer postgresPool.Close()

	if err = (SystemSetup{
		SeedAdminName:    cfg.Application.SeedAdminName,
		BackofficePRName: cfg.Application.BackofficePRName,
	}).Setup(context.Background(), postgresPool); err != nil {
		return fmt.Errorf("failed to run app setup: %w", err)
	}

	backofficePrId, openIDScopeId, err := resolveApplicationIDs(context.Background(), postgresPool, cfg.Application.BackofficePRName)
	if err != nil {
		return fmt.Errorf("failed to resolve application IDs: %w", err)
	}

	workerPool := workers.NewPool(10)
	outboxRepo := eventinfra.NewOutboxRepository(postgresPool)
	outboxEventBus := outboxapp.NewEventBus(outboxRepo, log)
	subsManager := modules.NewSubscriptionManager()
	cmdBus := commandInfra.NewCommandBus()
	qryBus := queryInfra.NewQueryBus()

	auditRepo := eventinfra.NewPostgresAuditRepository(postgresPool)
	subsManager.AddInMemorySubscription(auditapp.NewSubscriber(auditRepo))

	modules.RegisterBackoffice(postgresPool, cmdBus, qryBus, outboxEventBus, subsManager, openIDScopeId)
	modules.RegisterSecurity(postgresPool, cmdBus, qryBus, outboxEventBus, subsManager, cfg.Application.Name)
	modules.RegisterProtocol(postgresPool, cmdBus, qryBus, outboxEventBus, subsManager, cfg.Application.Name)

	eventBus := subsManager.Build(log, workerPool)

	registry := event.NewBuilderRegistry()
	modules.RegisterEventBuilders(registry)

	port, err := strconv.Atoi(cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("invalid server port: %w", err)
	}

	presenter := error_presenters.NewHttpErrorPresenter()
	managers := []app_errors.ErrorManager{error_managers.NewLogErrorPresenter(log)}
	hm := middlewares.NewHandlerMiddleware()

	ctx, srv := server.New(
		context.Background(),
		cfg.Server.Host,
		port,
		cfg.Server.ShutdownTimeout(),
		cmdBus,
		qryBus,
		log,
		presenter,
		managers,
		route_providers.NewBackoffice(hm, qryBus, cmdBus, postgresPool, backofficePrId, cfg.Application.Name),
		route_providers.NewProtocol(hm, cmdBus, qryBus, managers...),
	)

	consumer := eventinfra.NewOutboxConsumer(outboxRepo, eventBus, registry, 500*time.Millisecond, log)
	go consumer.Run(ctx)

	return srv.Run(ctx, log)
}
