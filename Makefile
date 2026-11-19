.PHONY: env up down migrate seed local-token dev-token run build test integration-test integration-test-up integration-test-down lint coverage
.PHONY: load-test stress-test

env:
	@test -n "$(ENV)" || (echo 'Usage: make env ENV=dev' >&2; exit 1)
	@bash ./scripts/create-env.sh "$(ENV)"

up:
	docker compose -f deploy/docker-compose.yml up -d --remove-orphans
	docker compose -f deploy/docker-compose.yml wait migrate

down:
	docker compose -f deploy/docker-compose.yml down -v --remove-orphans

migrate:
	@test -n "$(ENV)" || (echo 'Usage: make migrate ENV=dev' >&2; exit 1)
	go run ./scripts/migrate --env-file .env.$(ENV)

seed:
	@test -n "$(ENV)" || (echo 'Usage: make seed ENV=dev' >&2; exit 1)
	ENV="$(ENV)" ./scripts/seed.sh

run:
	ENV=$(or $(ENV),local) go run ./cmd/api

build:
	go build -o build/gas ./cmd/api

test:
	go test ./...

load-test:
	k6 run test/load/protocol.js

stress-test:
	k6 run -e PROFILE=stress test/load/protocol.js

integration-test-up:
	docker compose -f test/integrationtest/docker-compose.yml up -d --wait --wait-timeout 90

integration-test: integration-test-up
	MIGRATIONS_DSN='postgresql://gas_admin:integration-admin@127.0.0.1:55432/gas_backoffice?sslmode=disable' RUN_MIGRATIONS=true go run ./scripts/migrate
	TEST_POSTGRES_DSN='postgresql://gas-backoffice-user:integration-app@127.0.0.1:55432/gas_backoffice?sslmode=disable' go test -tags=integration -run '^TestIntegration' -count=1 -timeout=5m ./... ./internal/bounded_contexts/security/_shared/client_projection/infrastructure ./internal/bounded_contexts/protocol/_shared/infrastructure

integration-test-down:
	docker compose -f test/integrationtest/docker-compose.yml down --volumes --remove-orphans

coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out

lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...
