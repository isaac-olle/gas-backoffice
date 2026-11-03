.PHONY: env terraform-plan terraform-apply up down migrate seed local-token dev-token run build test lint coverage

env:
	@test -n "$(ENV)" || (echo 'Usage: make env ENV=dev' >&2; exit 1)
	@bash ./scripts/create-env.sh "$(ENV)"

terraform-plan:
	@bash ./scripts/terraform-plan.sh "$(ENV)" plan

terraform-apply:
	@bash ./scripts/terraform-plan.sh "$(ENV)" apply

up:
	docker compose -f deploy/docker-compose.yml up -d --remove-orphans
	docker compose -f deploy/docker-compose.yml wait migrate

down:
	docker compose -f deploy/docker-compose.yml down -v --remove-orphans

migrate:
	@test -n "$(ENV)" || (echo 'Usage: make migrate ENV=dev' >&2; exit 1)
	go run ./scripts/migrate --env-file .env.$(ENV)

seed:
	./scripts/seed.sh

run:
	ENV=$(or $(ENV),local) go run ./cmd/api

build:
	go build -o build/gas ./cmd/api

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out

lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...
