.PHONY: up down seed token run build test lint

up:
	docker compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down

seed:
	go run ./cmd/db_seed --env-file config/.env

token:
	go run ./scripts/generate_admin_token

run:
	go run ./cmd/api

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
