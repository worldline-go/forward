.DEFAULT_GOAL := help

.PHONY: run
run: ## Run the application
	go run cmd/forward/main.go

.PHONY: build
build: ## Build the binary
	goreleaser build --snapshot --clean --single-target

lint: ## Lint Go files
	golangci-lint --version
	golangci-lint run ./...

test: ## Run unit tests
	@go test -race ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
