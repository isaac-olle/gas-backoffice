# Integration tests

Integration files use the integration build tag and test names starting with
TestIntegration. The default command runs unit tests only:

    go test ./...

## Run locally

Requires Go 1.25.7 or newer and Docker Compose. From the repository root:

    make integration-test

This starts disposable PostgreSQL 17 on 127.0.0.1:55432, applies all migrations
with the deployment runner, and runs only integration tests. Tests connect as
gas-backoffice-user to check the application's database permissions too.
No development or production database or .env file is used.

Without Make, run these separate steps in PowerShell:

    docker compose -f test/integrationtest/docker-compose.yml up -d --wait --wait-timeout 90
    $env:MIGRATIONS_DSN = 'postgresql://gas_admin:integration-admin@127.0.0.1:55432/gas_backoffice?sslmode=disable'
    $env:RUN_MIGRATIONS = 'true'
    go run ./scripts/migrate
    # Continue only if migrations succeeded.
    $env:TEST_POSTGRES_DSN = 'postgresql://gas-backoffice-user:integration-app@127.0.0.1:55432/gas_backoffice?sslmode=disable'
    go test -tags=integration -run '^TestIntegration' -count=1 -timeout=5m ./... ./internal/bounded_contexts/security/_shared/client_projection/infrastructure ./internal/bounded_contexts/protocol/_shared/infrastructure

The explicit _shared paths include packages that Go's ./... pattern skips.
-count=1 prevents cached test results. Database tests fail if
TEST_POSTGRES_DSN is missing or the database cannot be reached.

Clean up afterwards, including after a failed test run:

    make integration-test-down

Without Make:

    docker compose -f test/integrationtest/docker-compose.yml down --volumes --remove-orphans

The database uses tmpfs and loses its data when stopped. The Compose project is
separate from deploy/docker-compose.yml.

## CI

The Integration tests job runs in parallel with Unit tests and Lint.
Separate steps start PostgreSQL, apply migrations, and run the Go integration
command on the runner. Failure prints PostgreSQL logs; cleanup always runs.
Building and deploying the application requires all three jobs to pass.

## Existing disposable database

Apply deploy/golang-migrate first, set TEST_POSTGRES_DSN, then run the tagged Go
command above. Tests create uniquely identified rows and clean them up when
finished. Never point integration tests at development or production.