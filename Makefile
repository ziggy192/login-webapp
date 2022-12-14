PWD = $(shell pwd)

SRC = `go list -f {{.Dir}} ./... | grep -v /vendor/`

all: fmt lint test

fmt:
	@echo "==> Formatting source code..."
	@go fmt $(SRC)

lint:
	@echo "==> Running lint check..."
	@go vet $(SRC)

test:
	@echo "==> Running tests..."
	@go clean -testcache ./...
	@go test `go list ./... | grep -v cmd` -race --cover

build-frontend:
	@go build -o frontend ./cmd/frontend

run-frontend:
	@./frontend

frontend: build-frontend run-frontend

build-api:
	@go build -o api ./cmd/api

run-api:
	@./api

api: build-api run-api

up:
	@docker compose \
		-f docker/docker-compose.dev.yml \
		-p ng_lu up -d

down:
	@docker compose \
		-f docker/docker-compose.dev.yml \
		-p ng_lu down -v --rmi local


.PHONY: all fmt lint test
