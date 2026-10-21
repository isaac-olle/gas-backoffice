package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	envFile := flag.String("env-file", "", "path to a .env.seed file to load (optional)")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil && !os.IsNotExist(err) {
			log.Fatalf("load env file %s: %v", *envFile, err)
		}
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "./deploy/golang-migrate"
	}
	migrationsDSN := os.Getenv("MIGRATIONS_DSN")
	if migrationsDSN == "" {
		log.Fatal("MIGRATIONS_DSN env var is required")
	}

	if os.Getenv("RUN_MIGRATIONS") == "false" {
		fmt.Println("migrations: skipped")
		return
	}

	m, err := migrate.New("file://"+migrationsPath, "pgx5://"+migrationsDSN[len("postgresql://"):])
	if err != nil {
		log.Fatalf("init migrations: %v", err)
	}
	defer func() { _, _ = m.Close() }()

	version, dirty, _ := m.Version()
	if dirty {
		log.Fatalf("database is dirty at version %d — fix manually:\n  psql \"$MIGRATIONS_DSN\" -c \"UPDATE schema_migrations SET dirty = false WHERE version = %d;\"", version, version)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("apply migrations: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("migrations: no changes")
	} else {
		fmt.Println("migrations: applied successfully")
	}
}
