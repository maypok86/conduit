PROJECT := "conduit"

BIN := "./bin/api"
SRC := "./cmd/api"

.PHONY: setup
setup: ## Install all the build and lint dependencies
	bash scripts/setup.sh

.PHONY: build
build: ## Build project
	bash scripts/build.sh $(BIN) $(SRC)

.PHONY: run
run: build ## Run project in local environment
	bash scripts/run.sh $(BIN)

.PHONY: up
up: ## Run project in docker environment
	bash scripts/up.sh $(PROJECT)

.PHONY: down
down: ## Stop project in docker environment
	docker compose -f deployments/docker-compose.yml -p $(PROJECT) --env-file .env down

.PHONY: logs
logs: ## View project logs from the docker container
	docker compose -f deployments/docker-compose.yml -p $(PROJECT) logs -f

.PHONY: fmt
fmt: ## Run format tools on all go files
	bash scripts/fmt.sh

.PHONY: lint
lint: ## Run all the linters
	golangci-lint run -v --color=always --timeout 4m ./...

.PHONY: test.unit
test.unit: ## Run all unit tests
	@echo 'mode: atomic' > coverage.txt
	go test -covermode=atomic -coverprofile=coverage.txt -v -race ./...

.PHONY: cover
cover: test.unit ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: ci
ci: lint test.unit ## Run all the tests and code checks

.PHONY: generate
generate: ## Generate files for the project
	go generate ./...
	mockgen --destination=internal/repository/psql/mocks/rows.go --package=mock_psql \
		--build_flags=--mod=mod github.com/jackc/pgx/v4 Rows
	mockgen --destination=internal/repository/psql/mocks/row.go --package=mock_psql \
		--build_flags=--mod=mod github.com/jackc/pgx/v4 Row
	mockgen --source=pkg/postgres/pgxpool.go --package=mock_psql \
		--destination=internal/repository/psql/mocks/pgxpool.go

.PHONY: clean
clean: ## Remove temporary files
	@go clean
	@rm -rf bin/
	@rm -rf coverage.txt
	@echo "SUCCESS!"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help
