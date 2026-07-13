package bootstrap

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/isaac-olle/gas-backoffice/internal/shared/logger"
)

func applyMigrations(log *logger.ZapLoggerAdapter, migrationsPath, migrationsDSN string) error {
	log.Info("running migrations")

	m, err := migrate.New("file://"+migrationsPath, "pgx5://"+migrationsDSN[len("postgresql://"):])
	if err != nil {
		return fmt.Errorf("failed to init migrations: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	version, dirty, _ := m.Version()
	if dirty {
		return fmt.Errorf("database is dirty at version %d — fix manually before starting the app", version)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("migrations: no changes")
	} else {
		log.Info("migrations: applied successfully")
	}

	return nil
}
