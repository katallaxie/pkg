.DEFAULT_GOAL := build

GO ?= go
GO_RUN_TOOLS ?= $(GO) run -modfile ./internal/tools/go.mod
GO_TEST = $(GO_RUN_TOOLS) gotest.tools/gotestsum --format pkgname

.PHONY: generate
generate: install ## Generate code.
	go generate ./...

.PHONY: install
install: ## Install tools.
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO_RUN_TOOLS) mvdan.cc/gofumpt -w .

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: generate fmt vet ## Run tests.
	mkdir -p .test/reports
	$(GO_TEST) --junitfile .test/reports/unit-test.xml -- -race ./... -count=1 -short -cover -coverprofile .test/reports/unit-test-coverage.out

.PHONY: lint
lint: ## Run golangci-lint against code.
	$(GO_RUN_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m -c .golangci.yml

.PHONY: clean
clean: ## Remove previous build.
	find . -type f -name '*.gen.go' -exec rm {} +
	git checkout go.mod
