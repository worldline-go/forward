BINARY  := forward
VERSION := $(or $(IMAGE_TAG),$(shell git describe --tags --first-parent --match "v*" 2> /dev/null || echo v0.0.0))
LOCAL_BIN_DIR := $(PWD)/bin

## golangci configuration
GOLANGCI_CONFIG_URL   := https://raw.githubusercontent.com/uber-go/guide/master/.golangci.yml
GOLANGCI_LINT_VERSION := v1.46.2

.PHONY: lint fmt test race msan coverage clean help $(BINARY) html html-gen html-wsl

$(BINARY): ## Build the binary file
	go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY) cmd/forward/main.go

.golangci.yml:
	curl -fkSL -o .golangci.yml $(GOLANGCI_CONFIG_URL)

bin/golangci-lint-$(GOLANGCI_LINT_VERSION):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN_DIR) $(GOLANGCI_LINT_VERSION)
	mv $(LOCAL_BIN_DIR)/golangci-lint $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)

lint: .golangci.yml bin/golangci-lint-$(GOLANGCI_LINT_VERSION) ## Lint Go files
	@$(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) --version
	@$(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) run ./...

fmt: ## Format Go files
	@go fmt -x ./...

test: ## Run unit tests
	@go test -race ./...

msan: ## Run memory sanitizer
	@go test -msan ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

html:
	@go tool cover -html=./coverage.out

html-gen: ## explorer.exe ./coverage.html
	@go tool cover -html=./coverage.out -o ./coverage.html

html-wsl: html-gen
	@explorer.exe `wslpath -w ./coverage.html`

clean: ## Remove previous build
	@rm -f $(BINARY)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
