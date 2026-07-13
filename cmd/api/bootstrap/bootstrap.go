package bootstrap

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	log.Info("running migrations")
	m, err := migrate.New("file://"+cfg.Databases.MigrationsPath, "pgx5://"+cfg.Databases.MigrationsDSN[len("postgresql://"):])
	if err != nil {
		return fmt.Errorf("failed to init migrations: %w", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	if err == migrate.ErrNoChange {
		log.Info("migrations: no changes")
	} else {
		log.Info("migrations: applied successfully")
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
	cmdBus := commandInfra.NewCommandBus()
	qryBus := queryInfra.NewQueryBus()

	auditRepo := eventinfra.NewPostgresAuditRepository(postgresPool)
	subsManager.AddInMemorySubscription(auditapp.NewSubscriber(auditRepo))

	modules.RegisterBackoffice(postgresPool, cmdBus, qryBus, outboxEventBus, subsManager)
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
		route_providers.NewBackoffice(hm, qryBus, cmdBus, postgresPool, cfg.Application.BackofficePrId, cfg.Application.Name),
		route_providers.NewProtocol(hm, cmdBus, qryBus),
	)

	consumer := eventinfra.NewOutboxConsumer(outboxRepo, eventBus, registry, 500*time.Millisecond, log)
	go consumer.Run(ctx)

	return srv.Run(ctx, log)
}
