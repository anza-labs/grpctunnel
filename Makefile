all: lint

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: tools
tools: ## Install all tools.
	go install tool

.PHONY: lint
lint: ## Run golangci-lint linters.
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint linters and perform fixes.
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --fix

.PHONY: fmt
fmt: ## Run golangci-lint formatters.
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint fmt

.PHONY: test
test: ## Run tests.
	go test -v -cover -coverprofile cover.out -race ./...

.PHONY: test-integration
test-integration: ## Run integration tests.
	go test -v -tags=integration -race ./internal/...

##@ CI

.PHONY: diff
diff: ## Run git diff-index to check if any changes are made.
	git --no-pager diff --exit-code HEAD --
