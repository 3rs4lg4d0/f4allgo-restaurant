# -----------------------------------------------------------------------------
# Database parameters for local development.
# -----------------------------------------------------------------------------
DB_HOST=localhost
DB_PORT=5432
DB_NAME=f4all_restaurant_service
DB_USER=f4all_restaurant_service
DB_PASSWORD=f4all_restaurant_pass

# -----------------------------------------------------------------------------
# Basic actions.
# -----------------------------------------------------------------------------
.PNOHY: all
all: generate build

.PHONY: generate
generate:
	@go generate ./...

.PHONY: build
build:
	@go build -o cmd/f4allgorestaurant-rest cmd/f4allgorestaurant-rest/f4allgorestaurant-rest.go
	@go build -o cmd/f4allgorestaurant-grpc cmd/f4allgorestaurant-grpc/f4allgorestaurant-grpc.go
	@go build -o cmd/f4allgorestaurant-cli cmd/f4allgorestaurant-cli/f4allgorestaurant-cli.go

.PHONY: clean-files
clean-files:
	@rm -f cmd/f4allgorestaurant-rest/f4allgorestaurant-rest
	@rm -f cmd/f4allgorestaurant-grpc/f4allgorestaurant-grpc
	@rm -f cmd/f4allgorestaurant-cli/f4allgorestaurant-cli
	@rm -f coverage.out

.PHONY: clean
clean: clean-files
	@go clean

.PHONY: full-clean
full-clean: clean-files
	@go clean -cache

.PHONY: test
test:
	@go test -cover ./...

.PHONY: test-coverage
test-coverage:
	@go test -v -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: run-rest-service
run-rest-service:
	@cd cmd/f4allgorestaurant-rest && go run f4allgorestaurant-rest.go

.PHONY: run-rest-service-debug
run-rest-service-debug:
	@export F4ALLGO_LOG_LEVEL=0 && export F4ALLGO_LOG_BEAUTIFY=true && \
		cd cmd/f4allgorestaurant-rest && go run f4allgorestaurant-rest.go

.PHONY: run-grpc-service
run-grpc-service:
	@cd cmd/f4allgorestaurant-grpc && go run f4allgorestaurant-grpc.go

.PHONY: mod-tidy
mod-tidy:
	@go mod tidy && cd tools && go mod tidy

# -----------------------------------------------------------------------------
# Action to migrate database. Use it when you need to test your SQL scripts.
# -----------------------------------------------------------------------------
.PHONY: migrate-db
migrate-db: command := up
migrate-db: db_url := "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
migrate-db:
	@echo "\n-----------------------------------------------------------------------------"
	@echo "> Executing database migrations"
	@echo "-----------------------------------------------------------------------------"
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -verbose -database $(db_url) -path sql $(command)

# -----------------------------------------------------------------------------
# Actions to manage local infrastructure using Docker containers.
# -----------------------------------------------------------------------------
.PHONY: start-infra
start-infra: script-start-local-infra migrate-db

.PHONY: stop-infra
stop-infra: script-stop-local-infra

.PHONY: delete-infra
delete-infra: script-delete-local-infra

# -----------------------------------------------------------------------------
# Utility targets (don't use them directly).
# -----------------------------------------------------------------------------
.PHONY: script-start-local-infra
script-start-local-infra:
	@echo "\n-----------------------------------------------------------------------------"
	@echo "> Starting docker containers"
	@echo "-----------------------------------------------------------------------------"
	@export POSTGRES_DB=$(DB_NAME) && \
	export POSTGRES_USER=$(DB_USER) && \
	export POSTGRES_PASSWORD=$(DB_PASSWORD) && \
	scripts/infra/start-local-infra.sh

.PHONY: script-start-local-infra
script-stop-local-infra:
	@echo "\n-----------------------------------------------------------------------------"
	@echo "> Stopping docker containers"
	@echo "-----------------------------------------------------------------------------"
	@export POSTGRES_DB=$(DB_NAME) && \
	export POSTGRES_USER=$(DB_USER) && \
	export POSTGRES_PASSWORD=$(DB_PASSWORD) && \
	scripts/infra/stop-local-infra.sh

.PHONY: script-delete-local-infra
script-delete-local-infra:
	@echo "\n-----------------------------------------------------------------------------"
	@echo "> Deleting docker containers"
	@echo "-----------------------------------------------------------------------------"
	@export POSTGRES_DB=$(DB_NAME) && \
	export POSTGRES_USER=$(DB_USER) && \
	export POSTGRES_PASSWORD=$(DB_PASSWORD) && \
	scripts/infra/delete-local-infra.sh