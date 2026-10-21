.PHONY: local-deploy local-undeploy up down migrate-local migrate-dev seed local-token dev-token run build test lint

up:
	docker compose -f deploy/docker-compose.yml up -d --remove-orphans
	docker compose -f deploy/docker-compose.yml wait migrate

down:
	docker compose -f deploy/docker-compose.yml down -v --remove-orphans

migrate-local:
	go run ./scripts/migrate --env-file .env.local

migrate-dev:
	go run ./scripts/migrate --env-file .env.dev


local-token:
	go run ./scripts/generate_admin_token --env .env.local

dev-token:
	go run ./scripts/generate_admin_token --env .env.dev


seed:
	./scripts/seed.sh


run:
	ENV=local go run ./cmd/api

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
