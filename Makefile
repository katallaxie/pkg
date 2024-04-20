.DEFAULT_GOAL := build

GO ?= go
GO_RUN_TOOLS 	?= $(GO) run -modfile ./tools/go.mod
GO_RELEASER 	?= $(GO_RUN_TOOLS) github.com/goreleaser/goreleaser
GO_TEST 		?= $(GO_RUN_TOOLS) gotest.tools/gotestsum --format pkgname

.PHONY: release
release: ## Release the project.
	$(GO_RELEASER) release --clean

.PHONY: generate
generate: install ## Generate code.
	$(GO) generate ./...

.PHONY: install
install: ## Install tools.
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO_RUN_TOOLS) mvdan.cc/gofumpt -w .

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: bench
bench: ## Run benchmarks.
	$(GO) test -bench=. --count=1 -run=^# ./...

.PHONY: test
test: fmt vet ## Run tests.
	mkdir -p .test/reports
	$(GO_TEST) --junitfile .test/reports/unit-test.xml -- -race ./... -count=1 -short -cover -coverprofile .test/reports/unit-test-coverage.out

.PHONY: lint
lint: ## Run golangci-lint against code.
	$(GO_RUN_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m -c .golangci.yml

.PHONY: clean
clean: ## Remove previous build.
	find . -type f -name '*.gen.go' -exec rm {} +
	git checkout go.mod
