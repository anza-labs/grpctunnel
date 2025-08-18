all: lint

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
	go test -cover -coverprofile cover.out -race ./...
