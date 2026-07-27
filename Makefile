.PHONY: up down migrate-local migrate-dev seed local-token dev-token run build test lint

up:
	docker compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down

migrate-local:
	go run ./cmd/migrate --env-file .env.local

migrate-dev:
	go run ./cmd/migrate --env-file .env.dev

seed:
	go run ./cmd/db_seed --env-file .env.dev

local-token:
	go run ./scripts/generate_admin_token --env .env.local

dev-token:
	go run ./scripts/generate_admin_token --env .env.dev

run:
	go run ./cmd/api --env-file .env.local

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
